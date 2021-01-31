package mecdsa

import (
	"errors"
	"io"
	"time"

	"bytes"
	"fmt"
	"math/big"

	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/log"
)

/*
DSMHandler 控制一次签名的完整周期
*/
type DSMHandler struct {
	Running bool //主线程启动控制

	db             *models.DB
	sessionID      common.Hash
	self           *models.NotaryInfo
	selfNotaryID   int
	otherNotaryIDs []int                  // 参与签名的其他人
	privateKeyInfo *models.PrivateKeyInfo //此次签名用到的分布式私钥
	notaryClient   notaryapi.NotaryClient

	xi                 share.SPrivKey            //协商好的私钥片
	paillierPubKeys    map[int]*proofs.PublicKey //其他公证人的同态加密公钥
	paillierPrivateKey *proofs.PrivateKey        //我的同态加密私钥
	publicKey          *share.SPubKey            //上次协商生成的总公钥
	vss                *feldman.VerifiableSS     //上次协商生成的feldman vss
	signMessage        *models.SignMessage       //此次签名生成过程中需要保存到数据库的信息

	phase2MessageBLock sync.Mutex //由于phase2Msg接收是在请求里面同步返回的,由于请求是并发的所以单独设置一个锁
	/*
		消息接收chan
	*/
	receiveChan chan api.Req
	/*
		流程控制chan
	*/
	phase1DoneChan    chan bool
	phase2DoneChan    chan bool
	phase3DoneChan    chan bool
	phase5A5BDoneChan chan bool
	phase5CDoneChan   chan bool
	phase6DoneChan    chan bool
	quitChan          chan error
}

/*
NewDSMHandler :
*/
func NewDSMHandler(db *models.DB, self *models.NotaryInfo, message messagetosign.MessageToSign, sessionID common.Hash, privateKeyInfo *models.PrivateKeyInfo, notaryClient notaryapi.NotaryClient) *DSMHandler {
	dh := &DSMHandler{
		db:                db,
		sessionID:         sessionID,
		self:              self,
		selfNotaryID:      self.ID,
		privateKeyInfo:    privateKeyInfo,
		notaryClient:      notaryClient,
		receiveChan:       make(chan api.Req, 10*cfg.Notaries.ShareCount),
		phase1DoneChan:    make(chan bool, 2),
		phase2DoneChan:    make(chan bool, 2),
		phase3DoneChan:    make(chan bool, 2),
		phase5A5BDoneChan: make(chan bool, 2),
		phase5CDoneChan:   make(chan bool, 2),
		phase6DoneChan:    make(chan bool, 2),
		quitChan:          make(chan error, 2),
	}
	signMessage := &models.SignMessage{
		Key:             dh.sessionID,
		UsedPrivateKey:  dh.privateKeyInfo.Key,
		Message:         message.GetSignBytes(),
		MessageName:     message.GetName(),
		SignTime:        time.Now().Unix(),
		Phase1BroadCast: make(map[int]*models.SignBroadcastPhase1),
		Phase2MessageB:  make(map[int]*models.MessageBPhase2),
		Phase3Delta:     make(map[int]*models.DeltaPhase3),
		Phase5A:         make(map[int]*models.Phase5A),
		Phase5C:         make(map[int]*models.Phase5C),
		Phase5D:         make(map[int]share.SPrivKey),
		AlphaGamma:      make(map[int]share.SPrivKey),
		AlphaWI:         make(map[int]share.SPrivKey),
		Delta:           make(map[int]share.SPrivKey),
	}
	dh.loadLockout(signMessage)
	return dh
}

// RegisterOtherNotaryIDs :
func (dh *DSMHandler) RegisterOtherNotaryIDs(otherNotaryIDs []int) *DSMHandler {
	var allNotaryIDs []int
	allNotaryIDs = append(allNotaryIDs, otherNotaryIDs...)
	allNotaryIDs = append(allNotaryIDs, dh.selfNotaryID)
	dh.otherNotaryIDs = otherNotaryIDs
	dh.signMessage.S = allNotaryIDs
	dh.createSignKeys()
	return dh
}

/*
StartDSMAndWaitFinish 主线程,负责控制一次dsm的流程
*/
func (dh *DSMHandler) StartDSMAndWaitFinish() (signature []byte, err error) {
	// 1. 人数校验
	if len(dh.signMessage.S) <= cfg.Notaries.ThresholdCount {
		err = fmt.Errorf("candidates notary too less")
		log.Info(sessionLogMsg(dh.sessionID, err.Error()))
		return
	}
	// 1.5 初始化phase1Msg
	phase1Msg := dh.generatePhase1Broadcast()
	// 2. 启动消息接收线程
	//go dh.receiveLoop()

	// 3. 开始phase1
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 1 start..."))
	phase1Req := notaryapi.NewDSMPhase1BroadcastRequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, phase1Msg)
	dh.notaryClient.WSBroadcast(phase1Req, dh.otherNotaryIDs...)
	// 等待phase1完成
	err = dh.waitPhaseDone(dh.phase1DoneChan)
	if err != nil {
		log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 1 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 1 done"))

	// 4. 开始phase2,phase2需要接收返回
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 2 start..."))
	phase2Msg := dh.generatePhase2MessageA()
	for _, notaryID := range dh.otherNotaryIDs {
		go func(notaryID int, phase2Msg *models.MessageA) {
			//defer wg.Done()
			phase2Req := notaryapi.NewDSMPhase2MessageARequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, phase2Msg)
			dh.notaryClient.SendWSReqToNotary(phase2Req, notaryID)
			resp, err2 := dh.notaryClient.WaitWSResponse(phase2Req.GetRequestID())
			if err2 != nil {
				err = err2
				log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 2 , err : %s ", err.Error()))
				dh.notify(nil, err)
				return
			}
			var respMsg models.MessageBPhase2
			err2 = resp.ParseData(&respMsg)
			if err2 != nil {
				err = err2
				log.Error("parse MessageBPhase2 err =%s \n%s", err.Error(), utils.ToJSONStringFormat(resp))
				dh.notify(nil, err)
				return
			}
			if respMsg.MessageBGamma == nil || respMsg.MessageBWi == nil {
				log.Error("parse MessageBPhase2 MessageBGamma nil   \n%s", utils.ToJSONStringFormat(resp))
				dh.notify(nil, errors.New("parse MessageBPhase2 MessageBGamma nil "))
				return
			}
			dh.receivePhase2MessageB(&respMsg, notaryID)
		}(notaryID, phase2Msg)
	}
	// 等待phase2完成
	err = dh.waitPhaseDone(dh.phase2DoneChan)
	if err != nil {
		log.Info(sessionLogMsg(dh.sessionID, "dsm failed at phase 2 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "phase 2 done"))

	// 5. 开始phase3
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 3 start..."))
	phase3Req := notaryapi.NewDSMPhase3DeltaIRequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, dh.generatePhase3DeltaI())
	dh.notaryClient.WSBroadcast(phase3Req, dh.otherNotaryIDs...)
	// 等待phase3完成
	err = dh.waitPhaseDone(dh.phase3DoneChan)
	if err != nil {
		log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 3 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 3 done"))

	// 6. 开始phase5A5B
	dh.generatePhase4R()
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 5A5B start..."))
	phase5A5BReq := notaryapi.NewDSMPhase5A5BProofRequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, dh.generatePhase5a5bZkProof())
	dh.notaryClient.WSBroadcast(phase5A5BReq, dh.otherNotaryIDs...)
	// 等待phase5A5B完成
	err = dh.waitPhaseDone(dh.phase5A5BDoneChan)
	if err != nil {
		log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 5A5B , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 5A5B done"))

	// 7. 开始phase5C
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 5C start..."))
	phase5CReq := notaryapi.NewDSMPhase5CProofRequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, dh.generatePhase5CProof())
	dh.notaryClient.WSBroadcast(phase5CReq, dh.otherNotaryIDs...)
	// 等待phase5C完成
	err = dh.waitPhaseDone(dh.phase5CDoneChan)
	if err != nil {
		log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 5C , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 5C done"))

	// 7. 开始phase6
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 6 start..."))
	phase6Req := notaryapi.NewDSMPhase6ReceiveSIRequest(dh.sessionID, dh.self, dh.privateKeyInfo.Key, dh.generate5dProof())
	dh.notaryClient.WSBroadcast(phase6Req, dh.otherNotaryIDs...)
	// 等待phase6完成
	err = dh.waitPhaseDone(dh.phase6DoneChan)
	if err != nil {
		log.Error(sessionLogMsg(dh.sessionID, "dsm failed at phase 6 , err : %s ", err.Error()))
		return
	}
	log.Info(sessionLogMsg(dh.sessionID, "dsm phase 6 done"))

	// 8.完成状态校验, 保存到DB并返回
	if dh.db != nil {
		err = dh.db.NewSignMessage(dh.signMessage)
		if err != nil {
			log.Error(sessionLogMsg(dh.sessionID, "save SignMessage to db err : %s ", err.Error()))
			return
		}
	}
	s := dh.signMessage.LocalSignature.SI.Clone()
	//所有人的的si，包括自己
	for i, si := range dh.signMessage.Phase5D {
		if i == dh.selfNotaryID {
			continue
		}
		share.ModAdd(s, si)
	}
	signature, verifyResult := verify(s, dh.signMessage.R, dh.signMessage.LocalSignature.Y, dh.signMessage.LocalSignature.M)
	if !verifyResult {
		err = errors.New("invilad signature")
	}
	return
}

/*
OnRequest 该次回话所有收到的消息的总入口
*/
func (dh *DSMHandler) OnRequest(req api.Req) {
	select {
	case dh.receiveChan <- req:
	default:
		log.Error("%T onRequest lost req=%s", dh, log.StringInterface(req, 3))
		// never block
	}
	//dh.receiveChan <- req
}

// 载入相关信息
func (dh *DSMHandler) loadLockout(signMessage *models.SignMessage) {
	p := dh.privateKeyInfo
	dh.xi = p.XI
	dh.paillierPrivateKey = p.PaillierPrivkey
	dh.paillierPubKeys = make(map[int]*proofs.PublicKey)
	for k, v := range p.PaillierKeysProof2 {
		dh.paillierPubKeys[k] = v.PaillierPubkey
	}
	dh.publicKey = &share.SPubKey{X: p.PublicKeyX, Y: p.PublicKeyY}
	dh.vss = p.SecretShareMessage3[dh.selfNotaryID].Vss

	dh.signMessage = signMessage
}

/*
//参数：自己的私钥片、系数点乘集合和我的多项式y集合、签名人的原始编号、所有签名人的编号
//==>每个公证人的公私钥、公钥片、{{{t,n},t+1个系数点乘G的结果(c1...c2)},y1...yn}
*/
func (dh *DSMHandler) createSignKeys() {
	lambdaI := dh.vss.MapShareToNewParams(dh.selfNotaryID, dh.signMessage.S) //lamda_i 解释：通过lamda_i对原所有证明人的群来映射出签名者群
	wi := share.ModMul(lambdaI, dh.xi)                                       //wi： 我原来的编号在对应签名群编号的映射关系 ，原来我是xi(私钥片) 现在是wi（我在新的签名群中的私钥片）
	gwiX, gwiY := share.S.ScalarBaseMult(wi.Bytes())                         //我在签名群中的公钥片
	gammaI := share.RandomPrivateKey()                                       //临时私钥
	gGammaIX, gGammaIY := share.S.ScalarBaseMult(gammaI.Bytes())             //临时公钥
	dh.signMessage.SignedKey = &models.SignedKey{
		WI:      wi,
		Gwi:     &share.SPubKey{X: gwiX, Y: gwiY},
		KI:      share.RandomPrivateKey(),
		GammaI:  gammaI,
		GGammaI: &share.SPubKey{X: gGammaIX, Y: gGammaIY},
	}
	return
}

/*
generatePhase1Broadcast 确定此次签名所用临时私钥,不能再换了.
*/
func (dh *DSMHandler) generatePhase1Broadcast() (msg *models.SignBroadcastPhase1) {
	blindFactor := share.RandomBigInt()
	gGammaIX, _ := share.S.ScalarBaseMult(dh.signMessage.SignedKey.GammaI.Bytes())
	com := createCommitmentWithUserDefinedRandomNess(gGammaIX, blindFactor)
	msg = &models.SignBroadcastPhase1{
		Com:         com,
		BlindFactor: blindFactor,
	}
	dh.signMessage.Phase1BroadCast[dh.selfNotaryID] = msg
	dh.checkPhase1Done()
	return
}

//receivePhase1Broadcast 收集此次签名中,其他公证人所用临时公钥,保证在后续步骤中不会被替换
func (dh *DSMHandler) receivePhase1Broadcast(msg *models.SignBroadcastPhase1, index int) {
	//dh.loadLockout(dh.signMessage)
	_, ok := dh.signMessage.Phase1BroadCast[index]
	if ok {
		dh.notify(nil, fmt.Errorf("pubkey roof for notary[ID=%d] already exist", index))
		return
	}
	dh.signMessage.Phase1BroadCast[index] = msg
	dh.checkPhase1Done()
	return
}

func (dh *DSMHandler) checkPhase1Done() {
	if len(dh.signMessage.Phase1BroadCast) == len(dh.signMessage.S) {
		dh.notify(dh.phase1DoneChan, nil)
	}
}

//p2p告知E(ki) 签名发起人，告诉其他所有人 步骤1 生成ki， 保证后续ki不会发生变化
func newMessageA(ki share.SPrivKey, paillierPubKey *proofs.PublicKey) (*models.MessageA, error) {
	ca, err := proofs.Encrypt(paillierPubKey, ki.Bytes())
	if err != nil {
		return nil, err
	}
	return &models.MessageA{C: ca}, nil
}

/*
GeneratePhase2MessageA p2p告知E(ki) 签名发起人，告诉其他所有人 步骤1 生成ki， 保证后续ki不会发生变化,
同时其他人是无法获取到ki,只有自己知道自己的ki
*/
func (dh *DSMHandler) generatePhase2MessageA() (msg *models.MessageA) {
	//dh.loadLockout(dh.signMessage)
	ma, err := newMessageA(dh.signMessage.SignedKey.KI, &dh.paillierPrivateKey.PublicKey)
	if err != nil {
		dh.notify(nil, err)
		return
	}
	dh.signMessage.MessageA = ma
	return ma
}

/*
NewMessageB :
我(bob)签名key中的临时私钥gammaI、(alice)paillier公钥、(alice)paillier公钥加密ki的结果
//返回：cB和两个证明(其他人验证)
我（alice）收到其他人（bob)给我的ma以后， 立即计算mb和两个证明，发送给bob
gammaI: alice的临时私钥
paillierPubKey: bob的同态加密公钥
ca: bob发送给alice的E(Ki)
*/
func newMessageB(gammaI share.SPrivKey, paillierPubKey *proofs.PublicKey, ca *models.MessageA) (*models.MessageB, error) {
	betaTagPrivateKey := share.RandomPrivateKey()
	cBetaTag, err := proofs.Encrypt(paillierPubKey, betaTagPrivateKey.Bytes()) //paillier加密一个随机数
	if err != nil {
		return nil, err
	}
	bca := proofs.Mul(paillierPubKey, ca.C, gammaI.Bytes()) //ca.C：加密ki的结果 gammaI:gammaI
	//cB=b * E(ca) + E(beta_tag)   (b:gammaI  ca:ca.C   )
	cb := proofs.AddCipher(paillierPubKey, bca, cBetaTag)
	//beta= -bata_tag mod q
	beta := share.ModSub(share.PrivKeyZero.Clone(), betaTagPrivateKey)

	//todo 提供证明 ：证明gammaI是我自己的,证明beta是我合法提供的 既然提供了beta,这里面的betatagproof应该是可以忽略的
	bproof := proofs.Prove(gammaI)
	betaTagProof := proofs.Prove(betaTagPrivateKey)
	return &models.MessageB{
		C:            cb,
		BProof:       bproof,
		BetaTagProof: betaTagProof,
		Beta:         beta,
	}, nil
}

/*
ReceivePhase2MessageA alice 收到来自bob的E(ki), 向bob提供证明,自己持有着gammaI和WI
其中gammaI是临时私钥
WI包含着XI
*/
func (dh *DSMHandler) receivePhase2MessageA(msg *models.MessageA, index int) (mb *models.MessageBPhase2) {
	//dh.loadLockout(dh.signMessage)
	mgGamma, err := newMessageB(dh.signMessage.SignedKey.GammaI, dh.paillierPubKeys[index], msg)
	if err != nil {
		panic(err)
	}
	mbw, err := newMessageB(dh.signMessage.SignedKey.WI, dh.paillierPubKeys[index], msg)
	if err != nil {
		panic(err)
	}
	mb = &models.MessageBPhase2{
		MessageBGamma: mgGamma,
		MessageBWi:    mbw,
	}
	return
}

func verifyProofsGetAlpha(m *models.MessageB, dk *proofs.PrivateKey, a share.SPrivKey) (share.SPrivKey, error) {
	ashare, err := proofs.Decrypt(dk, m.C) //用dk解密cB
	if err != nil {
		return share.SPrivKey{}, err
	}
	alpha := new(big.Int).SetBytes(ashare)
	alphaKey := share.BigInt2PrivateKey(alpha)
	gAlphaX, gAlphaY := share.S.ScalarBaseMult(alphaKey.Bytes())
	babTagX, babTagY := share.S.ScalarMult(m.BProof.PK.X, m.BProof.PK.Y, a.Bytes())
	babTagX, babTagY = share.PointAdd(babTagX, babTagY, m.BetaTagProof.PK.X, m.BetaTagProof.PK.Y)
	if proofs.Verify(m.BProof) && proofs.Verify(m.BetaTagProof) &&
		babTagX.Cmp(gAlphaX) == 0 &&
		babTagY.Cmp(gAlphaY) == 0 {
		return alphaKey, nil
	}
	return share.SPrivKey{}, errors.New("invalid key")
}

//receivePhase2MessageB 收集并验证MessageB证明信息
func (dh *DSMHandler) receivePhase2MessageB(msg *models.MessageBPhase2, index int) {
	dh.phase2MessageBLock.Lock()
	dh.signMessage.Phase2MessageB[index] = msg
	dh.phase2MessageBLock.Unlock()
	dh.checkPhase2Done()
	return
}

func (dh *DSMHandler) checkPhase2Done() {
	if len(dh.signMessage.Phase2MessageB) == len(dh.signMessage.S)-1 {
		for index, msg := range dh.signMessage.Phase2MessageB {
			alphaijGamma, err := verifyProofsGetAlpha(msg.MessageBGamma, dh.paillierPrivateKey, dh.signMessage.SignedKey.KI)
			if err != nil {
				dh.notify(nil, err)
				return
			}
			alphaijWi, err := verifyProofsGetAlpha(msg.MessageBWi, dh.paillierPrivateKey, dh.signMessage.SignedKey.KI)
			if err != nil {
				dh.notify(nil, err)
				return
			}
			dh.signMessage.AlphaGamma[index] = alphaijGamma
			dh.signMessage.AlphaWI[index] = alphaijWi
		}
		dh.notify(dh.phase2DoneChan, nil)
	}
}

func (dh *DSMHandler) phase2DeltaI() share.SPrivKey {
	var k *models.SignedKey
	k = dh.signMessage.SignedKey
	if len(dh.signMessage.AlphaGamma) != len(dh.signMessage.S)-1 {
		panic("arg error")
	}
	//kiGammaI=ki * gammI+Sum(alpha_vec) +Sum(beta_vec)
	kiGammaI := k.KI.Clone()
	share.ModMul(kiGammaI, k.GammaI)
	for _, i := range dh.signMessage.S {
		if i == dh.selfNotaryID {
			continue
		}
		share.ModAdd(kiGammaI, dh.signMessage.AlphaGamma[i])
		share.ModAdd(kiGammaI, dh.signMessage.Phase2MessageB[i].MessageBGamma.Beta)
	}
	return kiGammaI
}
func (dh *DSMHandler) phase2SigmaI() share.SPrivKey {
	if len(dh.signMessage.AlphaWI) != len(dh.signMessage.S)-1 {
		panic("length error")
	}
	kiwi := dh.signMessage.SignedKey.KI.Clone()
	share.ModMul(kiwi, dh.signMessage.SignedKey.WI)
	//todo vij=vji ?
	for _, i := range dh.signMessage.S {
		if i == dh.selfNotaryID {
			continue
		}
		share.ModAdd(kiwi, dh.signMessage.AlphaWI[i])
		share.ModAdd(kiwi, dh.signMessage.Phase2MessageB[i].MessageBWi.Beta)
	}
	return kiwi
}

/*
generatePhase3DeltaI 依据上一步协商信息,生成我自己的DeltaI,然后广播给所有其他人,需要这些参与公证人得到完整的的Delta
但是生成的SigmaI自己保留,在生成自己的签名片的时候使用
*/
func (dh *DSMHandler) generatePhase3DeltaI() (msg *models.DeltaPhase3) {
	//dh.loadLockout(dh.signMessage)
	deltaI := dh.phase2DeltaI()
	sigmaI := dh.phase2SigmaI()
	dh.signMessage.Sigma = sigmaI
	dh.signMessage.Delta[dh.selfNotaryID] = deltaI
	msg = &models.DeltaPhase3{Delta: deltaI}
	dh.checkPhase3Done()
	return
}

//receivePhase3DeltaI 收集所有的deltaI
func (dh *DSMHandler) receivePhase3DeltaI(msg *models.DeltaPhase3, index int) {
	//dh.loadLockout(dh.signMessage)
	_, ok := dh.signMessage.Delta[index]
	if ok {
		dh.notify(nil, fmt.Errorf("ReceivePhase3DeltaI for %d already exist", index))
		return
	}
	dh.signMessage.Delta[index] = msg.Delta
	dh.checkPhase3Done()
	return
}

func (dh *DSMHandler) checkPhase3Done() {
	if len(dh.signMessage.Delta) == len(dh.signMessage.S) {
		dh.notify(dh.phase3DoneChan, nil)
	}
}

func phase4(delta share.SPrivKey,
	gammaIProof map[int]*models.MessageBPhase2,
	phase1Broadcast map[int]*models.SignBroadcastPhase1) (*share.SPubKey, error) {
	if len(gammaIProof) != len(phase1Broadcast) {
		panic("length must equal")
	}
	for i, p := range gammaIProof {
		//校验第一步广播的参数没有发生变化
		if createCommitmentWithUserDefinedRandomNess(p.MessageBGamma.BProof.PK.X,
			phase1Broadcast[i].BlindFactor).Cmp(phase1Broadcast[i].Com) == 0 {
			continue
		}
		return nil, errors.New("invliad key")
	}

	//tao_i=g^^gamma
	sumx, sumy := new(big.Int), new(big.Int)
	for _, p := range gammaIProof {
		if sumx.Cmp(big.NewInt(0)) == 0 && sumy.Cmp(big.NewInt(0)) == 0 {
			sumx = p.MessageBGamma.BProof.PK.X
			sumy = p.MessageBGamma.BProof.PK.Y
		} else {
			sumx, sumy = share.PointAdd(sumx, sumy, p.MessageBGamma.BProof.PK.X, p.MessageBGamma.BProof.PK.Y)
		}

	}
	rx, ry := share.S.ScalarMult(sumx, sumy, delta.Bytes())
	return &share.SPubKey{X: rx, Y: ry}, nil
}

//phase3 计算：inverse(delta) mod q
// all parties broadcast delta_i and compute delta_i ^(-1)
func phase3ReconstructDelta(delta map[int]share.SPrivKey) share.SPrivKey {
	sum := share.PrivKeyZero.Clone()
	for _, deltaI := range delta {
		share.ModAdd(sum, deltaI)
	}
	return share.InvertN(sum)
}

//generatePhase4R 所有公证人都应该得到相同的R,其中R.X就是最后签名(r,s,v)中的r
func (dh *DSMHandler) generatePhase4R() (R *share.SPubKey) {
	//dh.loadLockout(dh.signMessage)
	delta := phase3ReconstructDelta(dh.signMessage.Delta)
	//少一个自己的MessageB todo fixme 这里面需要有更好的方式来生成数据,后续必须优化
	mgGamma, err := newMessageB(dh.signMessage.SignedKey.GammaI, &dh.paillierPrivateKey.PublicKey, dh.signMessage.MessageA)
	if err != nil {
		dh.notify(nil, err)
		return
	}
	dh.signMessage.Phase2MessageB[dh.selfNotaryID] = &models.MessageBPhase2{MessageBGamma: mgGamma, MessageBWi: nil}
	R, err = phase4(delta, dh.signMessage.Phase2MessageB, dh.signMessage.Phase1BroadCast)
	if err != nil {
		return
	}
	dh.signMessage.R = R
	return
}

//const
func phase5LocalSignature(ki share.SPrivKey, message *big.Int,
	R *share.SPubKey, sigmaI share.SPrivKey,
	pubkey *share.SPubKey) *models.LocalSignature {
	m := share.BigInt2PrivateKey(message)
	r := share.BigInt2PrivateKey(R.X)
	si := share.ModMul(m, ki)
	share.ModMul(r, sigmaI)
	share.ModAdd(si, r) //si=m * k_i + r * sigma_i
	return &models.LocalSignature{
		LI:   share.RandomPrivateKey(),
		RhoI: share.RandomPrivateKey(),
		//li:   big.NewInt(71),
		//rhoi: big.NewInt(73),
		R: &share.SPubKey{
			X: new(big.Int).Set(R.X),
			Y: new(big.Int).Set(R.Y),
		},
		SI: si, //签名片
		M:  new(big.Int).Set(message),
		Y: &share.SPubKey{
			X: new(big.Int).Set(pubkey.X),
			Y: new(big.Int).Set(pubkey.Y),
		},
	}

}

//com:commit的Ci(广播)
//Phase5ADecom1 :commit的Di(广播)
func phase5aBroadcast5bZkproof(dh *models.LocalSignature) (*models.Phase5Com1, *models.Phase5ADecom1, *proofs.HomoELGamalProof) {
	blindFactor := share.RandomBigInt()
	//Ai=g^^rho_i
	aix, aiy := share.S.ScalarBaseMult(dh.RhoI.Bytes())
	lIRhoI := dh.LI.Clone()
	share.ModMul(lIRhoI, dh.RhoI)
	//Bi=G*lIRhoI
	bix, biy := share.S.ScalarBaseMult(lIRhoI.Bytes())
	//vi=R*si+G*li
	tx, ty := share.S.ScalarMult(dh.R.X, dh.R.Y, dh.SI.Bytes()) //R^^si
	vix, viy := share.S.ScalarBaseMult(dh.LI.Bytes())           //g^^li
	vix, viy = share.PointAdd(vix, viy, tx, ty)

	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{
		{X: vix, Y: viy}, {X: aix, Y: aiy}, {X: bix, Y: biy},
	})
	com := createCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)

	//proof是5b的zkp构造
	witness := proofs.NewHomoElGamalWitness(dh.LI, dh.SI) //li si
	delta := &proofs.HomoElGamalStatement{
		G: share.NewGE(aix, aiy),               //Ai
		H: share.NewGE(dh.R.X, dh.R.Y),         //R
		Y: share.NewGE(share.S.Gx, share.S.Gy), //g
		D: share.NewGE(vix, viy),               //Vi
		E: share.NewGE(bix, biy),               //Bi
	}
	//证明提供的是正确的si???
	proof := proofs.CreateHomoELGamalProof(witness, delta)
	return &models.Phase5Com1{Com: com},
		&models.Phase5ADecom1{
			Vi:          share.NewGE(vix, viy),
			Ai:          share.NewGE(aix, aiy),
			Bi:          share.NewGE(bix, biy),
			BlindFactor: blindFactor,
		},
		proof
}

/*
generatePhase5a5bZkProof 从此步骤开始,互相不断交换信息来判断对方能够生成正确的Si(也就是签名片),
如果所有参与者都能生成最终的签名片,那么我才能把自己的签名片告诉对方.
si的累加和就是签名(r,s,v)中的s
*/
func (dh *DSMHandler) generatePhase5a5bZkProof() (msg *models.Phase5A) {
	//dh.loadLockout(dh.signMessage)
	//messageHash := utils.Sha256(dh.signMessage.Message)
	messageBN := new(big.Int).SetBytes(dh.signMessage.Message[:])
	localSignature := phase5LocalSignature(dh.signMessage.SignedKey.KI, messageBN, dh.signMessage.R, dh.signMessage.Sigma, dh.publicKey)
	phase5Com, phase5ADecom, helgamalProof := phase5aBroadcast5bZkproof(localSignature)
	msg = &models.Phase5A{Phase5Com1: phase5Com, Phase5ADecom1: phase5ADecom, Proof: helgamalProof}
	dh.signMessage.Phase5A[dh.selfNotaryID] = msg
	dh.signMessage.LocalSignature = localSignature
	dh.checkPhase5A5BDone()
	return
}

// receivePhase5A5BProof :
func (dh *DSMHandler) receivePhase5A5BProof(msg *models.Phase5A, index int) {
	//dh.loadLockout(dh.signMessage)
	_, ok := dh.signMessage.Phase5A[index]
	if ok {
		dh.notify(nil, fmt.Errorf("ReceivePhase5A5BProof already exist for %d", index))
		return
	}
	dh.signMessage.Phase5A[index] = msg
	dh.checkPhase5A5BDone()
	return
}

func (dh *DSMHandler) checkPhase5A5BDone() {
	if len(dh.signMessage.Phase5A) == len(dh.signMessage.S) {
		for _, msg := range dh.signMessage.Phase5A {
			delta := &proofs.HomoElGamalStatement{
				G: msg.Phase5ADecom1.Ai,
				H: dh.signMessage.R,
				Y: share.NewGE(share.S.Gx, share.S.Gy),
				D: msg.Phase5ADecom1.Vi,
				E: msg.Phase5ADecom1.Bi,
			}
			inputhash := proofs.CreateHashFromGE([]*share.SPubKey{
				msg.Phase5ADecom1.Vi,
				msg.Phase5ADecom1.Ai,
				msg.Phase5ADecom1.Bi,
			})
			e := createCommitmentWithUserDefinedRandomNess(inputhash.D, msg.Phase5ADecom1.BlindFactor)
			if e.Cmp(msg.Phase5Com1.Com) == 0 &&
				msg.Proof.Verify(delta) {

			} else {
				dh.notify(nil, errors.New("invalid com"))
				return
			}
		}
		dh.notify(dh.phase5A5BDoneChan, nil)
	}
}

func phase5c(dh *models.LocalSignature, deCommitments []*models.Phase5ADecom1,
	vi *share.SPubKey) (*models.Phase5Com2, *models.Phase5DDecom2) {

	//从广播的commit(Ci,Di)得到vi,ai
	v := vi.Clone()
	for i := 0; i < len(deCommitments); i++ {
		v.X, v.Y = share.PointAdd(v.X, v.Y, deCommitments[i].Vi.X, deCommitments[i].Vi.Y)
	}
	a := deCommitments[0].Ai.Clone()
	for i := 1; i < len(deCommitments); i++ {
		a.X, a.Y = share.PointAdd(a.X, a.Y, deCommitments[i].Ai.X, deCommitments[i].Ai.Y)
	}
	r := share.BigInt2PrivateKey(dh.R.X)
	yrx, yry := share.S.ScalarMult(dh.Y.X, dh.Y.Y, r.Bytes())
	m := share.BigInt2PrivateKey(dh.M)
	//Vi之积×g^(-m)*y^(-r)
	gmx, gmy := share.S.ScalarBaseMult(m.Bytes())
	v.X, v.Y = share.PointSub(v.X, v.Y, gmx, gmy)
	v.X, v.Y = share.PointSub(v.X, v.Y, yrx, yry)
	//UI=V * rhoi
	uix, uiy := share.S.ScalarMult(v.X, v.Y, dh.RhoI.Bytes())
	//Ti=A * li
	tix, tiy := share.S.ScalarMult(a.X, a.Y, dh.LI.Bytes())

	//commit(UI ,Ti)，广播出去
	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{
		{X: uix, Y: uiy},
		{X: tix, Y: tiy},
	})
	blindFactor := share.RandomBigInt()
	com := createCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)
	return &models.Phase5Com2{Com: com},
		&models.Phase5DDecom2{
			UI:          &share.SPubKey{X: uix, Y: uiy}, //Ci
			Ti:          &share.SPubKey{X: tix, Y: tiy}, //Di
			BlindFactor: blindFactor,
		}
}

//generatePhase5CProof  fixme 提供一个好的注释
func (dh *DSMHandler) generatePhase5CProof() (msg *models.Phase5C) {
	//dh.loadLockout(dh.signMessage)
	if len(dh.signMessage.Phase5A) != len(dh.signMessage.S) {
		panic("cannot genrate5c until all 5b proof received")
	}
	var decomVec []*models.Phase5ADecom1
	for i, m := range dh.signMessage.Phase5A {
		if i == dh.selfNotaryID { //phase5c 不应该包括自己的decommitment
			continue
		}
		decomVec = append(decomVec, m.Phase5ADecom1)
	}
	phase5com2, phase5decom2 := phase5c(dh.signMessage.LocalSignature, decomVec, dh.signMessage.Phase5A[dh.selfNotaryID].Phase5ADecom1.Vi)
	msg = &models.Phase5C{Phase5Com2: phase5com2, Phase5DDecom2: phase5decom2}
	dh.signMessage.Phase5C[dh.selfNotaryID] = msg
	dh.checkPhase5CDone()
	return
}

//receivePhase5cProof   fixme 暂时没有好的解释
func (dh *DSMHandler) receivePhase5cProof(msg *models.Phase5C, index int) {
	//dh.loadLockout(dh.signMessage)
	//验证5c的hash(ui,ti)=5c的Ci

	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{msg.Phase5DDecom2.UI, msg.Phase5DDecom2.Ti})
	inputhash.D = createCommitmentWithUserDefinedRandomNess(inputhash.D, msg.Phase5DDecom2.BlindFactor)
	if inputhash.D.Cmp(msg.Phase5Com2.Com) != 0 {
		dh.notify(nil, errors.New("invalid com"))
		return
	}
	_, ok := dh.signMessage.Phase5C[index]
	if ok {
		dh.notify(nil, fmt.Errorf("ReceivePhase5cProof for %d already exist", index))
		return
	}
	dh.signMessage.Phase5C[index] = msg
	dh.checkPhase5CDone()
	return
}

func (dh *DSMHandler) checkPhase5CDone() {
	if len(dh.signMessage.Phase5C) == len(dh.signMessage.S) {
		dh.notify(dh.phase5CDoneChan, nil)
	}
}

/*
generate5dProof 接受所有签名人的si的广播，有可能某个公证人会保留信息，最终生成有效的签名，私自保留下来,但是不告诉其他人自己的si是多少.
但是这种情况其他公证人可以知道,没有收到某个公证人的si
*/
func (dh *DSMHandler) generate5dProof() (si share.SPrivKey) {
	//dh.loadLockout(dh.signMessage)
	var deCommitments2 []*models.Phase5DDecom2
	var commitments2 []*models.Phase5Com2
	var deCommitments1 []*models.Phase5ADecom1
	for i, m := range dh.signMessage.Phase5C {
		deCommitments2 = append(deCommitments2, m.Phase5DDecom2)
		commitments2 = append(commitments2, m.Phase5Com2)
		deCommitments1 = append(deCommitments1, dh.signMessage.Phase5A[i].Phase5ADecom1)
	}
	//校验收到的关于si的信息,都是真实的,正确的,只有通过了,才能把自己的si告诉给其他公证人
	si, err := phase5d(dh.signMessage.LocalSignature, deCommitments2, commitments2, deCommitments1)
	if err != nil {
		dh.notify(nil, err)
		return
	}
	dh.signMessage.Phase5D[dh.selfNotaryID] = si
	dh.checkPhase6Done()
	return
}

/*
 校验收到的关于si的信息,都是真实的,正确的,只有通过了,才能把自己的si告诉给其他公证人
*/
func phase5d(dh *models.LocalSignature, deCommitments2 []*models.Phase5DDecom2,
	commitments2 []*models.Phase5Com2,
	deCommitments1 []*models.Phase5ADecom1) (share.SPrivKey, error) {
	if len(deCommitments1) != len(deCommitments2) ||
		len(deCommitments2) != len(commitments2) {
		panic("arg error")
	}

	biasedSumTbX := new(big.Int).Set(share.S.Gx)
	biasedSumTbY := new(big.Int).Set(share.S.Gy)

	for i := 0; i < len(commitments2); i++ {
		//(5c的ti + 5a的bi)连加
		biasedSumTbX, biasedSumTbY = share.PointAdd(biasedSumTbX, biasedSumTbY,
			deCommitments2[i].Ti.X, deCommitments2[i].Ti.Y)
		biasedSumTbX, biasedSumTbY = share.PointAdd(biasedSumTbX, biasedSumTbY,
			deCommitments1[i].Bi.X, deCommitments1[i].Bi.Y)
	}
	//用于比较 UI 和 (5c的ti + 5a的bi)连加 是否相等
	for i := 0; i < len(commitments2); i++ {
		biasedSumTbX, biasedSumTbY = share.PointSub(
			biasedSumTbX, biasedSumTbY,
			deCommitments2[i].UI.X, deCommitments2[i].UI.Y,
		)
	}
	//log.Trace(fmt.Sprintf("(gx,gy)=(%s,%s)", share.S.Gx.Text(16), share.S.Gy.Text(16)))
	//log.Trace(fmt.Sprintf("(tbx,tby)=(%s,%s)", biasedSumTbX.Text(16), biasedSumTbY.Text(16)))
	if share.S.Gx.Cmp(biasedSumTbX) == 0 &&
		share.S.Gy.Cmp(biasedSumTbY) == 0 {
		return dh.SI.Clone(), nil
	}
	return share.PrivKeyZero, errors.New("invalid key")
}

//recevieSI 收集签名片,收集齐以后就可以得到完整的签名.所有公证人都应该得到有效的签名.
func (dh *DSMHandler) recevieSI(si share.SPrivKey, index int) {
	//dh.loadLockout(dh.signMessage)
	if _, ok := dh.signMessage.Phase5D[index]; ok {
		dh.notify(nil, fmt.Errorf("si for %d already received", index))
		return
	}
	dh.signMessage.Phase5D[index] = si
	dh.checkPhase6Done()
	return
}

func (dh *DSMHandler) checkPhase6Done() {
	if len(dh.signMessage.Phase5D) == len(dh.signMessage.S) {
		s := dh.signMessage.LocalSignature.SI.Clone()
		//所有人的的si，包括自己
		for i, si := range dh.signMessage.Phase5D {
			if i == dh.selfNotaryID {
				continue
			}
			share.ModAdd(s, si)
		}
		_, verifyResult := verify(s, dh.signMessage.R, dh.signMessage.LocalSignature.Y, dh.signMessage.LocalSignature.M)
		if !verifyResult {
			dh.notify(nil, errors.New("invilad signature"))
		} else {
			dh.notify(dh.phase6DoneChan, nil)
		}
	}
}

/*
使用SignatureNormalize来对签名进行处理,符合EIP155签名要求.
https://ethereum.stackexchange.com/questions/42455/during-ecdsa-signing-how-do-i-generate-the-recovery-id
I never found any proper documentation about the Recovery ID but I did talk with somebody on Reddit and they gave me my answer:

id = y1 & 1; // Where (x1,y1) = k x G;
if (s > curve.n / 2) id = id ^ 1; // Invert id if s of signature is over half the n
I had to modify the mbedtls library to pass back the Recovery ID but when I did I could generate transactions that Geth accepted 100% of the time.

The long explanation:

During signing, a point is generated (X, Y) called R and a number called S. R's X goes on to become r and S becomes s. In order to generate the Recovery ID you take the one's bit from Y. If S is bigger than half the curve's N parameter you invert that bit. That bit is the Recovery ID. Ethereum goes on to manipulate it to indicate compressed or uncompressed addresses as well as indicate what chain the transaction was signed for (so the transaction can't be replayed on another Ethereum chain that the private key might be present on). These modifications to the Recovery ID become v.

There's also a super rare chance that you need to set the second bit of the recovery id meaning the recovery id could in theory be 0, 1, 2, or 3. But there's a 0.000000000000000000000000000000000000373% of needing to set the second bit according to a question on Bitcoin.SE.


*/
func verify(s share.SPrivKey, R, y *share.SPubKey, message *big.Int) ([]byte, bool) {
	r := share.BigInt2PrivateKey(R.X)
	b := share.InvertN(s)
	a := share.BigInt2PrivateKey(message)
	u1 := a.Clone()
	u1 = share.ModMul(u1, b)
	u2 := r.Clone()
	u2 = share.ModMul(u2, b)

	gu1x, gu1y := share.S.ScalarBaseMult(u1.Bytes())
	yu2x, yu2y := share.S.ScalarMult(y.X, y.Y, u2.Bytes())
	gu1x, gu1y = share.PointAdd(gu1x, gu1y, yu2x, yu2y)
	//已经确认是一个有效的签名
	if share.BigInt2PrivateKey(gu1x).D.Cmp(r.D) == 0 {
		//return true
	} else {
		return nil, false
	}
	sig := make([]byte, 65)
	copy(sig[:32], bigIntTo32Bytes(r.D))
	copy(sig[32:64], bigIntTo32Bytes(s.D))
	//leave sig[64]=0

	//按照以太坊EIP155的规范来处理签名
	h := common.Hash{}
	var err error
	h.SetBytes(message.Bytes())
	if s.D.Cmp(halfN) >= 0 {
		/* //所谓的normalize就是如果s>n/2,s=n-s ,保证s的唯一性
		sig, err = secp256k1.SignatureNormalize(sig)
		if err != nil {
			log.Error(fmt.Sprintf("SignatureNormalize err %s\n,r=%s,s=%s", err, r.D, s.D))
			return nil, false
		}*/
		s2 := new(big.Int)
		s2 = s2.Sub(share.S.N, s.D)
		copy(sig[32:64], bigIntTo32Bytes(s2))
		tmp := readBigInt(bytes.NewBuffer(sig[32:64]))
		tmp = tmp.Add(tmp, s.D)
		if tmp.Cmp(share.S.N) != 0 {
			panic("must equal")
		}
		//log.Info(fmt.Sprintf("s=%s,normals=%s,n=fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", s.D.Text(16), .Text(16)))
	}

	key, err := crypto.GenerateKey()
	pubkey := key.PublicKey
	pubkey.X = y.X
	pubkey.Y = y.Y
	addr := crypto.PubkeyToAddress(pubkey)

	/*
		id = y1 & 1; // Where (x1,y1) = k x G; x1,y1 is R
		if (s > curve.n / 2) id = id ^ 1; // Invert id if s of signature is over half the n
	*/
	v := new(big.Int).Mod(R.Y, big2).Int64()
	if s.D.Cmp(halfN) > 0 {
		v = v ^ 1
	}
	sig[64] = byte(v)
	//try v=0
	pubkeybin, err := crypto.Ecrecover(h[:], sig)
	if err == nil {
		pubkey2 := crypto.ToECDSAPub(pubkeybin)
		addr2 := crypto.PubkeyToAddress(*pubkey2)
		if addr2 == addr {
			restoreAndCheckRS(sig)
			return sig, true
		}

	} else {
		log.Error(fmt.Sprintf("Ecrecover err %s", err))
	}
	return nil, false
}

//bigIntTo32Bytes convert a big int to bytes
func bigIntTo32Bytes(i *big.Int) []byte {
	data := i.Bytes()
	buf := make([]byte, 32)
	for i := 0; i < 32-len(data); i++ {
		buf[i] = 0
	}
	for i := 32 - len(data); i < 32; i++ {
		buf[i] = data[i-32+len(data)]
	}
	return buf
}

//readBigInt read big.Int from buffer
func readBigInt(reader io.Reader) *big.Int {
	bi := new(big.Int)
	tmpbuf := make([]byte, 32)
	_, err := reader.Read(tmpbuf)
	if err != nil {
		log.Error(fmt.Sprintf("read BigInt error %s", err))
	}
	bi.SetBytes(tmpbuf)
	return bi
}

func restoreAndCheckRS(sig []byte) {
	buf := bytes.NewBuffer(sig)
	readBigInt(buf)
	s := readBigInt(buf)
	v, err := buf.ReadByte()
	if err != nil {
		panic(err)
	}
	if v != 0 && v != 1 {
		panic("wrong v")
	}

	if s.Cmp(halfN) >= 0 {
		panic("wrong s")
	}
}

var halfN *big.Int
var big2 *big.Int

func init() {
	halfN = new(big.Int).Set(share.S.N)
	big2 = big.NewInt(2)
	halfN.Div(halfN, big2)
}

func (dh *DSMHandler) waitPhaseDone(c chan bool) (err error) {
	var req api.Req
	for {
		select {
		case <-c:
			return
		case err = <-dh.quitChan:
			if err != nil {
				log.Error(sessionLogMsg(dh.sessionID, "waitPhaseDone of DSMHandler quit with err %s", err.Error()))
			}
			return
		case req = <-dh.receiveChan:
			switch r := req.(type) {
			case *notaryapi.DSMPhase1BroadcastRequest:
				dh.receivePhase1Broadcast(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.DSMPhase2MessageARequest:
				resp := dh.receivePhase2MessageA(r.Msg, r.GetSenderNotaryID())
				r.WriteSuccessResponse(resp)
			case *notaryapi.DSMPhase3DeltaIRequest:
				dh.receivePhase3DeltaI(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.DSMPhase5A5BProofRequest:
				dh.receivePhase5A5BProof(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.DSMPhase5CProofRequest:
				dh.receivePhase5cProof(r.Msg, r.GetSenderNotaryID())
			case *notaryapi.DSMPhase6ReceiveSIRequest:
				dh.recevieSI(r.Msg, r.GetSenderNotaryID())
			default:
				log.Error(sessionLogMsg(dh.sessionID, "unknown msg for DSMHandler :\n%s", utils.ToJSONStringFormat(req)))
			}
		}
	}
}

func (dh *DSMHandler) notify(c chan bool, err error) {
	if err != nil {
		// 这里写两遍.因为可能有两个线程
		select {
		case dh.quitChan <- err:
		default:
			log.Error("%T notify err  lost,err=%s", dh, err)
			// never block
		}
		select {
		case dh.quitChan <- err:
		default:
			log.Error("%T notify err  lost,err=%s", dh, err)
			// never block
		}
	} else {
		//c <- true
		select {
		case c <- true:
		default:
			panic("write to channel blocked")
			// never block
		}
	}
}
