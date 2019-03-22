package mecdsa

import (
	"fmt"
	"math/big"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
PKNHandler 控制一次私钥协商的完整周期
*/
type PKNHandler struct {
	db             *models.DB
	sessionID      common.Hash
	self           *models.NotaryInfo
	selfNotaryID   int
	otherNotaryIDs []int
	privateKeyInfo *models.PrivateKeyInfo
	notaryapi.NotaryClient

	/*
		消息接收chan
	*/
	receiveChan chan api.Req
	/*
		流程控制chan
	*/
	phase1DoneChan chan bool
	phase2DoneChan chan bool
	phase3DoneChan chan bool
	phase4DoneChan chan bool
	quitChan       chan error
}

/*
NewPKNHandler :
*/
func NewPKNHandler(db *models.DB, self *models.NotaryInfo, otherNotaryIDs []int, sessionID common.Hash, notaryClient notaryapi.NotaryClient) *PKNHandler {
	ui, dk := createKeys()
	privateKeyInfo := &models.PrivateKeyInfo{
		Key:                 sessionID,
		UI:                  ui,
		PaillierPrivkey:     dk,
		Status:              models.PrivateKeyNegotiateStatusInit,
		PubKeysProof1:       make(map[int]*models.KeyGenBroadcastMessage1),
		PaillierKeysProof2:  make(map[int]*models.KeyGenBroadcastMessage2),
		SecretShareMessage3: make(map[int]*models.KeyGenBroadcastMessage3),
		LastPubkeyProof4:    make(map[int]*models.KeyGenBroadcastMessage4),
		CreateTime:          time.Now().Unix(),
	}
	return &PKNHandler{
		db:             db,
		sessionID:      sessionID,
		self:           self,
		selfNotaryID:   self.ID,
		otherNotaryIDs: otherNotaryIDs,
		privateKeyInfo: privateKeyInfo,
		receiveChan:    make(chan api.Req, 10*params.ShareCount),
		phase1DoneChan: make(chan bool, 2),
		phase2DoneChan: make(chan bool, 2),
		phase3DoneChan: make(chan bool, 2),
		phase4DoneChan: make(chan bool, 2),
		quitChan:       make(chan error, 2),
		NotaryClient:   notaryClient,
	}
}

/*
StartPKNAndWaitFinish 开始一次私钥协商
主动发起时,参数req传空
*/
func (ph *PKNHandler) StartPKNAndWaitFinish(req *notaryapi.KeyGenerationPhase1MessageRequest) (privateKeyInfo *models.PrivateKeyInfo, err error) {
	// 0. 启动消息处理线程
	//go ph.receiveLoop()
	// 0.5 投递
	if req != nil {
		ph.OnRequest(req)
	}
	// ======> 1. 开始phase1
	log.Info(sessionLogMsg(ph.sessionID, "phase 1 start..."))
	phase1Req := notaryapi.NewKeyGenerationPhase1MessageRequest(ph.sessionID, ph.self, ph.generatePhase1PubKeyProof())
	ph.WSBroadcast(phase1Req, ph.otherNotaryIDs...)
	// 等待phase1完成
	err = ph.waitPhaseDone(ph.phase1DoneChan)
	if err != nil {
		log.Info(sessionLogMsg(ph.sessionID, "pkn failed at phase 1 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(ph.sessionID, "phase 1 done"))

	// 2. 开始phase2
	log.Info(sessionLogMsg(ph.sessionID, "phase 2 start..."))
	phase2Req := notaryapi.NewKeyGenerationPhase2MessageRequest(ph.sessionID, ph.self, ph.generatePhase2PaillierKeyProof())
	ph.WSBroadcast(phase2Req, ph.otherNotaryIDs...)
	// 等待phase2完成
	err = ph.waitPhaseDone(ph.phase2DoneChan)
	if err != nil {
		log.Info(sessionLogMsg(ph.sessionID, "pkn failed at phase 2 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(ph.sessionID, "phase 2 done"))

	//  3. 开始phase3
	log.Info(sessionLogMsg(ph.sessionID, "phase 3 start..."))
	phase3MsgMap := ph.generatePhase3SecretShare()
	for notaryID, phase3Msg := range phase3MsgMap {
		// 按ID分别发送phase3消息给其他人
		phase3Req := notaryapi.NewKeyGenerationPhase3MessageRequest(ph.sessionID, ph.self, phase3Msg)
		ph.SendWSReqToNotary(phase3Req, notaryID)
	}
	// 等待phase3完成
	err = ph.waitPhaseDone(ph.phase3DoneChan)
	if err != nil {
		log.Info(sessionLogMsg(ph.sessionID, "pkn failed at phase 3 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(ph.sessionID, "phase 3 done"))

	// 4. 开始phase4
	log.Info(sessionLogMsg(ph.sessionID, "phase 4 start..."))
	phase4Req := notaryapi.NewKeyGenerationPhase4MessageRequest(ph.sessionID, ph.self, ph.generatePhase4PubKeyProof())
	ph.WSBroadcast(phase4Req, ph.otherNotaryIDs...)
	// 等待phase4完成
	err = ph.waitPhaseDone(ph.phase4DoneChan)
	if err != nil {
		log.Info(sessionLogMsg(ph.sessionID, "pkn failed at phase 4 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(ph.sessionID, "phase 4 done and private key for %s ready", ph.privateKeyInfo.ToAddress().String()))
	// 5.完成状态校验, 保存到DB并返回
	ph.privateKeyInfo.Address = ph.privateKeyInfo.ToAddress()
	if ph.db != nil {
		err = ph.db.NewPrivateKeyInfo(ph.privateKeyInfo)
		if err != nil {
			log.Error(sessionLogMsg(ph.sessionID, "save PrivateKeyInfo to db err : %s ", err.Error()))
			return
		}
	}
	privateKeyInfo = ph.privateKeyInfo
	return
}

// OnRequest 消息入口
func (ph *PKNHandler) OnRequest(req api.Req) {
	select {
	case ph.receiveChan <- req:
	default:
		// never block
	}
}

/*
generatePhase1PubKeyProof phase1.1 生成自己的随机数,这是后续feldman vss的种子
同时也生成了第二步所需的同态加密公钥私钥.
这一步生成的同态公钥需要在第二步告诉其他所有公证人
这一步生成的随机数对应的公钥片,需要告诉所有其他公证人.
*/
func (ph *PKNHandler) generatePhase1PubKeyProof() (msg *models.KeyGenBroadcastMessage1) {
	//用dlogproof 传播我自己的公钥片,包含相关的零知识证明
	proof := proofs.Prove(ph.privateKeyInfo.UI)
	msg = &models.KeyGenBroadcastMessage1{Proof: proof}
	return
}

/*
receivePhase1PubKeyProof phase 1.2 接受其他公证人传递过来的公钥片证明信息,
如果凑齐了所有公证人(params.ShareCount)的公钥片信息,那么就可以组合出来此次协商最终的公钥.
但是到目前为止没有一个公证人知道最终公钥对应的私钥片
*/
func (ph *PKNHandler) receivePhase1PubKeyProof(m *models.KeyGenBroadcastMessage1, index int) {
	//fmt.Printf("notary[ID=%d] receive phase 1 msg from notary[ID=%d]\n", ph.selfNotaryID, index)
	p := ph.privateKeyInfo
	if !proofs.Verify(m.Proof) {
		ph.notify(nil, fmt.Errorf("pubkey roof for notary[ID=%d] verify not pass", index))
		return
	}
	_, ok := p.PubKeysProof1[index]
	if ok {
		ph.notify(nil, fmt.Errorf("pubkey roof for notary[ID=%d] already exist", index))
		return
	}
	p.PubKeysProof1[index] = m
	ph.checkPhase1Done()
	return
}

func (ph *PKNHandler) checkPhase1Done() {
	p := ph.privateKeyInfo
	//除了自己以外的因素都凑齐了
	if len(p.PubKeysProof1) == params.ShareCount-1 {
		x, y := share.S.ScalarBaseMult(p.UI.Bytes())
		for _, m := range p.PubKeysProof1 {
			x, y = share.PointAdd(x, y, m.Proof.PK.X, m.Proof.PK.Y)
		}
		//所有人公证人都掌握的总公钥
		p.PublicKeyX = x
		p.PublicKeyY = y
		ph.notify(ph.phase1DoneChan, nil)
	}
}

/*
generatePhase2PaillierKeyProof phase 2.1 广播给其他所有公证人的同态加密公钥以及Proof
*/
func (ph *PKNHandler) generatePhase2PaillierKeyProof() (msg *models.KeyGenBroadcastMessage2) {
	p := ph.privateKeyInfo
	blindFactor := share.RandomBigInt()
	//生成关于同态加密私钥的零知识证明
	correctKeyProof := proofs.CreateNICorrectKeyProof(p.PaillierPrivkey)
	x, _ := share.S.ScalarBaseMult(p.UI.Bytes())
	com := createCommitmentWithUserDefinedRandomNess(x, blindFactor)
	msg = &models.KeyGenBroadcastMessage2{
		Com:             com,
		PaillierPubkey:  &p.PaillierPrivkey.PublicKey,
		CorrectKeyProof: correctKeyProof,
		BlindFactor:     blindFactor,
	}
	p.PaillierKeysProof2[ph.selfNotaryID] = msg
	ph.checkPhase2Done()
	return
}

//receivePhase2PaillierPubKeyProof phase 2.2 收到其他公证人的同态加密公钥信息
func (ph *PKNHandler) receivePhase2PaillierPubKeyProof(m *models.KeyGenBroadcastMessage2, index int) {
	//fmt.Printf("notary[ID=%d] receive phase 2 msg from notary[ID=%d]\n", ph.selfNotaryID, index)
	//对方证明私钥我确实知道,虽然没有告诉我
	if !m.CorrectKeyProof.Verify(m.PaillierPubkey) {
		ph.notify(nil, fmt.Errorf("paillier pubkey roof for notary[ID=%d] verify not pass", index))
		return
	}
	p := ph.privateKeyInfo
	//不接受重复的消息
	_, ok := p.PaillierKeysProof2[index]
	if ok {
		ph.notify(nil, fmt.Errorf("paillier pubkey roof for notary[ID=%d] already exist", index))
		return
	}
	p.PaillierKeysProof2[index] = m
	ph.checkPhase2Done()
	return
}

func (ph *PKNHandler) checkPhase2Done() {
	p := ph.privateKeyInfo
	//所有因素都凑齐了
	if len(p.PaillierKeysProof2) == params.ShareCount {
		ph.notify(ph.phase2DoneChan, nil)
	}
}

/*
generatePhase3SecretShare phase 3.1基于第一步的随机数生成SecretShares,
假设有三个公证人,我的编号是0
0 生成的shares[s00,s01,s02]
1 生成的shares[s10,s11,s12]
2 生成的shares[s20,s21,s22]
那么分发过程是
0: s01->1 s02->2
1: s10->0 s12->2
2: s20->0 s21->1
最终收集齐以后:
0 持有[s00,s10,s20] 累加即为私钥片x0
1 持有[s01,s11,s21] 累加即为私钥片x1
2 持有[202,s12,s22] 累加即为私钥片x2
每个公证人务必保管好自己的xi,这是后续签名必须.
*/
func (ph *PKNHandler) generatePhase3SecretShare() (msgs map[int]*models.KeyGenBroadcastMessage3) {
	p := ph.privateKeyInfo
	for i := 0; i < len(p.PubKeysProof1); i++ {
		if i == ph.selfNotaryID {
			continue
		}
		m := p.PaillierKeysProof2[i]
		if createCommitmentWithUserDefinedRandomNess(p.PubKeysProof1[i].Proof.PK.X, m.BlindFactor).Cmp(m.Com) != 0 {
			ph.notify(nil, fmt.Errorf("blind factor error for notary[ID=%d]", i))
			return
		}
	}
	vss, secretShares := feldman.Share(params.ThresholdCount, params.ShareCount, p.UI)
	msg := &models.KeyGenBroadcastMessage3{
		Vss:         vss,
		SecretShare: secretShares[ph.selfNotaryID],
		Index:       ph.selfNotaryID,
	}
	p.SecretShareMessage3[ph.selfNotaryID] = msg
	msgs = make(map[int]*models.KeyGenBroadcastMessage3)
	for i := 0; i < params.ShareCount; i++ {
		if i == ph.selfNotaryID {
			continue
		}
		msg := &models.KeyGenBroadcastMessage3{
			Vss:         vss,
			SecretShare: secretShares[i],
			Index:       i,
		}
		msgs[i] = msg
	}
	ph.checkPhase3Done()
	return
}

//receivePhase3SecretShare 接收来自其他公证人定向分发给我的secret share.
func (ph *PKNHandler) receivePhase3SecretShare(msg *models.KeyGenBroadcastMessage3, index int) {
	//fmt.Printf("notary[ID=%d] receive phase 3 msg from notary[ID=%d]\n", ph.selfNotaryID, index)
	p := ph.privateKeyInfo
	//vss的第0号证据对应的就是第一步随机数对应的公钥片,必须用的是那一个
	if !equalGE(p.PubKeysProof1[index].Proof.PK, msg.Vss.Commitments[0]) {
		ph.notify(nil, fmt.Errorf("pubkey not match for notary[ID=%d]", index))
		return
	}
	_, ok := p.SecretShareMessage3[index]
	if ok {
		ph.notify(nil, fmt.Errorf("secret shares for notary[ID=%d] already received", index))
		return
	}
	//必须经过feldman vss验证,符合规则.如果某个公证人给的secret share发错了,比如s11发送给了2号公证人,这里会检测出错误.
	if !msg.Vss.ValidateShare(msg.SecretShare, ph.selfNotaryID+1) {
		ph.notify(nil, fmt.Errorf("secret share error for notary[ID=%d]", index))
		return
	}
	p.SecretShareMessage3[index] = msg
	ph.checkPhase3Done()
	return
}

func (ph *PKNHandler) checkPhase3Done() {
	p := ph.privateKeyInfo
	//所有因素都凑齐了.可以计算我自己持有的私钥片 sum(f(i)) f(i)解释:公证人生成了n个多项式 ,把所有第i个多项式的结果相加)
	if len(p.SecretShareMessage3) == params.ShareCount {
		p.XI = share.SPrivKey{D: new(big.Int)}
		for _, s := range p.SecretShareMessage3 {
			share.ModAdd(p.XI, s.SecretShare)
		}
		ph.notify(ph.phase3DoneChan, nil)
	}
}

/*
generatePhase4PubKeyProof phase4.1 生成我自己对我持有私钥片的证明,
todo 目前问题就在于对方无法验证对方给我的XI是真实有效的,和公钥是关联在一起的真实私钥片,
有可能对方在GeneratePhase3SecretShare 就广播一个错误的secret share,但是我也无法证明
*/
func (ph *PKNHandler) generatePhase4PubKeyProof() (msg *models.KeyGenBroadcastMessage4) {
	p := ph.privateKeyInfo
	if p.XI.D == nil {
		panic("must wait for phase 3 finish")
	}
	proof := proofs.Prove(p.XI)
	msg = &models.KeyGenBroadcastMessage4{Proof: proof}
	// 保存
	p.LastPubkeyProof4[ph.selfNotaryID] = msg

	ph.checkPhase4Done()
	return
}

/*
receivePhase4VerifyTotalPubKey phase4.2 接收对方持有的Xi的证明,并验证其有效性. 但是我怎么知道他是有效的呢? 目前为止是没有办法的.
*/
func (ph *PKNHandler) receivePhase4VerifyTotalPubKey(msg *models.KeyGenBroadcastMessage4, index int) {
	p := ph.privateKeyInfo
	if !proofs.Verify(msg.Proof) {
		ph.notify(nil, fmt.Errorf("last pubkey verify error for notary[ID=%d]", index))
		return
	}
	_, ok := p.LastPubkeyProof4[index]
	if ok {
		ph.notify(nil, fmt.Errorf("last pubkey for notary[ID=%d] already exist", index))
		return
	}
	p.LastPubkeyProof4[index] = msg
	ph.checkPhase4Done()
	return
}

func (ph *PKNHandler) checkPhase4Done() {
	p := ph.privateKeyInfo
	if len(p.LastPubkeyProof4) == params.ShareCount {
		//可以校验pubkey之和是否有效
		x, y := new(big.Int), new(big.Int)
		x.Set(p.LastPubkeyProof4[0].Proof.PK.X)
		y.Set(p.LastPubkeyProof4[0].Proof.PK.Y)
		for i := 1; i < len(p.LastPubkeyProof4); i++ {
			x, y = share.PointAdd(x, y, p.LastPubkeyProof4[i].Proof.PK.X, p.LastPubkeyProof4[i].Proof.PK.Y)
		}
		if x.Cmp(p.PublicKeyX) != 0 || y.Cmp(p.PublicKeyY) != 0 {
			//panic("should equal")
		}
		p.Status = models.PrivateKeyNegotiateStatusFinished
		ph.notify(ph.phase4DoneChan, nil)
	}
}
func (ph *PKNHandler) waitPhaseDone(c chan bool) (err error) {
	var req api.Req
	for {
		select {
		case <-c:
			return
		case err = <-ph.quitChan:
			if err != nil {
				log.Error(sessionLogMsg(ph.sessionID, "waitPhaseDone of PKNHandler quit with err %s", err.Error()))
			}
		case req = <-ph.receiveChan:
			switch r := req.(type) {
			case *notaryapi.KeyGenerationPhase1MessageRequest:
				ph.receivePhase1PubKeyProof(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.KeyGenerationPhase2MessageRequest:
				ph.receivePhase2PaillierPubKeyProof(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.KeyGenerationPhase3MessageRequest:
				ph.receivePhase3SecretShare(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.KeyGenerationPhase4MessageRequest:
				ph.receivePhase4VerifyTotalPubKey(r.Msg, r.GetSenderNotaryID())
			default:
				log.Error(sessionLogMsg(ph.sessionID, "unknown msg for PKNHandler :\n%s", utils.ToJSONStringFormat(req)))
			}
		}
	}
}

func (ph *PKNHandler) notify(c chan bool, err error) {
	if err != nil {
		// 这里写两遍.因为可能有两个线程
		select {
		case ph.quitChan <- err:
		default:
			// never block
		}
		select {
		case ph.quitChan <- err:
		default:
			// never block
		}
	} else {
		//c <- true
		select {
		case c <- true:
		default:
			// never block
		}
	}
}
