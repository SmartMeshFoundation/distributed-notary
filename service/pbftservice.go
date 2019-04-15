package service

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/pbft/pbft"
	utils "github.com/nkbai/goutils"
	"github.com/nkbai/log"
)

/*
负责节点之间的nonce协商
*/
type pbftService struct {
	key       string //协商哪一个地址的nonce
	client    *pbft.Client
	clientMsg chan interface{}

	server    *pbft.Server
	serverMsg chan interface{}

	dispatchService dispatchServiceBackend
	allNotaries     []*models.NotaryInfo
	notaryClient    notaryapi.NotaryClient
	nonces          map[string]chan pbft.OpResult
	lock            sync.RWMutex
	db              *models.DB
	quit            chan struct{}
}

func NewPBFTService(key string, allNotaries []*models.NotaryInfo, notaryClient notaryapi.NotaryClient, dispatchService dispatchServiceBackend, db *models.DB) *pbftService {
	ps := &pbftService{
		clientMsg:       make(chan interface{}, 10),
		serverMsg:       make(chan interface{}, 10),
		notaryClient:    notaryClient,
		dispatchService: dispatchService,
		key:             key,
		nonces:          make(map[string]chan pbft.OpResult),
		db:              db,
		quit:            make(chan struct{}),
	}
	myid := dispatchService.getSelfNotaryInfo().ID
	var allids []int
	for i := 0; i < len(allNotaries); i++ {
		allids = append(allids, allNotaries[i].ID)
	}
	f := len(allNotaries) / 3
	nonce, err := db.GetNonce(key)
	if err != nil {
		panic(err)
	}
	nonce++ //数据库中存的是上次用的nonce,
	log.Trace(fmt.Sprintf("allids=%v,nonce=%d", allids, nonce))
	ps.client = pbft.NewPBFTClient(myid, ps.clientMsg, ps, f, allids)
	ps.server = pbft.NewPBFTServer(myid, f, nonce, ps.serverMsg, ps, allids, ps)
	go ps.loop()
	return ps
}

//SendMessage 这里是否应该处理一下
func (ps *pbftService) SendMessage(msg interface{}, target int) {
	req := &notaryapi.PBFTMessage{
		BaseReq:              api.NewBaseReq(notaryapi.APINamePBFTMessage),
		BaseReqWithSignature: api.NewBaseReqWithSignature(ps.dispatchService.getSelfNotaryInfo().GetAddress()),
		Key:                  ps.key,
		Msg:                  pbft.EncodeMsg(msg),
	}
	//log.Trace(fmt.Sprintf("ps sendMessage %v,to %d", msg, target))
	if target == ps.dispatchService.getSelfNotaryInfo().ID {
		ps.OnRequest(req)
	} else {
		ps.notaryClient.SendWSReqToNotary(req, target)
	}
}

//OnRequest 来自其他公证人的pbft消息
func (ps *pbftService) OnRequest(req *notaryapi.PBFTMessage) {
	var n *models.NotaryInfo

	if req.GetSigner() != ps.dispatchService.getSelfNotaryInfo().GetAddress() {
		var ok bool
		n, ok = ps.dispatchService.getNotaryService().getNotaryInfoByAddress(req.GetSigner())
		if !ok {
			log.Error(fmt.Sprintf("receive req,but signer is unkown,req=%s", utils.StringInterface(req, 3)))
			return
		}
	} else {
		n = ps.dispatchService.getSelfNotaryInfo()
	}

	/*
		todo 需要解决公证人id问题,需要带上签名,
	*/
	msg := pbft.DecodeMsg(req.Msg)
	switch msg2 := msg.(type) {
	case pbft.ClientMessager:
		if !msg2.CheckMessageSender(n.ID) {
			log.Error(fmt.Sprintf("receive pbft message,but signer not match,signer=%d,msg=%s",
				n.ID, utils.StringInterface(msg, 3)))
			return
		}
		ps.clientMsg <- msg
	case pbft.ServerMessager:
		if !msg2.CheckMessageSender(n.ID) {
			log.Error(fmt.Sprintf("receive pbft message,but signer not match,signer=%d,msg=%s",
				n.ID, utils.StringInterface(msg, 3)))
			return
		}
		ps.serverMsg <- msg
	default:
		log.Error(fmt.Sprintf("pbftService onRquest unkown req=%s", utils.StringInterface(req, 2)))
	}

}

func (ps *pbftService) Stop() {
	close(ps.quit)
}
func (ps *pbftService) loop() {
	for {
		select {
		case op := <-ps.client.Apply:
			ps.lock.RLock()
			c, ok := ps.nonces[op.Cmd]
			ps.lock.RUnlock()
			if !ok {
				log.Warn(fmt.Sprintf("%s receive pbft notify,but have no related channel", ps.key))
				continue //ignore
			}
			c <- op
		case <-ps.quit:
			return
		}

	}
}

func (ps *pbftService) newNonce(op string) (nonce uint64, err error) {
	log.Trace(fmt.Sprintf("ps[%s] new nonce for %s", ps.key, op))
	ps.lock.Lock()
	c, ok := ps.nonces[op]
	if ok {
		ps.lock.Unlock()
		err = fmt.Errorf("already exist req %s", op)
		return
	}
	c = make(chan pbft.OpResult, 1)
	ps.nonces[op] = c
	ps.lock.Unlock()
	ps.client.Start(op)
	r := <-c
	log.Trace(fmt.Sprintf("ps[%s] new nonce return %s", ps.key, utils.StringInterface(r, 3)))
	return uint64(r.Seq), r.Error
}

func (ps *pbftService) UpdateSeq(seq int, _, _ string) {
	nonce, err := ps.db.GetNonce(ps.key)
	if err != nil {
		panic(err)
	}
	if seq > nonce {
		err = ps.db.UpdateNonce(ps.key, seq)
		if err != nil {
			panic(err)
		}
	}
}

/*
	GetOpAuxiliary 根据来自用户的op构造相应的辅助信息,
	对于以太坊来说,就很简单,就是op的hash值
	对于比特币来说就是,分配出去的UTXO列表
*/
func (ps *pbftService) GetOpAuxiliary(op string, _ int) string {
	return pbft.Digest(op)
}

/*
	VerifyAuxiliary 在集齐验证prepare消息后,验证op对应的auxiliary是否有效.
	对于以太坊来说,总是有效的
	对于比特币来说,可能因为分配出去的utxo已经使用,金额不够等原因造成失败
*/
func (ps *pbftService) VerifyAuxiliary(_ int, _ string, _ string) error {
	return nil
}
