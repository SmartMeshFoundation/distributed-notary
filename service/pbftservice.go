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
PBFTService 负责节点之间的nonce协商
*/
type PBFTService struct {
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

// NewPBFTService :
func NewPBFTService(key string, allNotaries []*models.NotaryInfo, notaryClient notaryapi.NotaryClient, dispatchService dispatchServiceBackend, db *models.DB) *PBFTService {
	ps := &PBFTService{
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
func (ps *PBFTService) SendMessage(msg interface{}, target int) {
	req := &notaryapi.PBFTMessage{
		BaseReq: api.BaseReq{
			Name: notaryapi.APINamePBFTMessage,
		},
		Key: ps.key,
		Msg: pbft.EncodeMsg(msg),
	}
	//log.Trace(fmt.Sprintf("ps sendMessage %v,to %d", msg, target))
	if target == ps.dispatchService.getSelfNotaryInfo().ID {
		ps.OnRequest(req)
	} else {
		ps.notaryClient.SendWSReqToNotary(req, target)
	}
}

// OnRequest :
func (ps *PBFTService) OnRequest(req *notaryapi.PBFTMessage) {
	msg := pbft.DecodeMsg(req.Msg)

	switch msg.(type) {
	case pbft.ClientMessager:
		ps.clientMsg <- msg
	case pbft.ServerMessager:
		ps.serverMsg <- msg
	default:
		log.Error(fmt.Sprintf("PBFTService onRquest unkown req=%s", utils.StringInterface(req, 2)))
	}

}

// Stop :
func (ps *PBFTService) Stop() {
	close(ps.quit)
}
func (ps *PBFTService) loop() {
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

func (ps *PBFTService) newNonce(op string) (nonce uint64, err error) {
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

// UpdateSeq :
func (ps *PBFTService) UpdateSeq(seq int) {
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
