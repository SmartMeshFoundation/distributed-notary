package mecdsa

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/goutils"
)

/*
ThresholdPrivKeyGenerator 一次协商私钥的过程,
所有参与的公证人都参与,最终形成一个分布式私钥,每个公证人只能持有一部分私钥片,
没有一个公证人知道完整的私钥. 但是后续这些公证人可以在不暴露自己私钥片的基础上进行签名.
*/
type ThresholdPrivKeyGenerator struct {
	selfNotaryID int
	db           *models.DB
	PrivateKeyID common.Hash //此次协商唯一的key
}

/*
NewThresholdPrivKeyGenerator 生成此次协商起始所需参数
key: 此次协商唯一标志
暗含的其他公证人都已知的信息,包括 ThresholdCount和ShareCount
*/
func NewThresholdPrivKeyGenerator(selfNotaryID int, db *models.DB, privateKeyID common.Hash) *ThresholdPrivKeyGenerator {
	return &ThresholdPrivKeyGenerator{
		selfNotaryID: selfNotaryID,
		db:           db,
		PrivateKeyID: privateKeyID,
	}
}

//phase1.1 生成自己的随机数,所有公证人的私钥片都会从这里面取走一部分
func createKeys() (share.SPrivKey, *proofs.PrivateKey) {
	ui := share.RandomPrivateKey()
	dk, err := proofs.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return ui, dk
}

/*
GeneratePhase1PubKeyProof phase1.1 生成自己的随机数,这是后续feldman vss的种子
同时也生成了第二步所需的同态加密公钥私钥.
这一步生成的同态公钥需要在第二步告诉其他所有公证人
这一步生成的随机数对应的公钥片,需要告诉所有其他公证人.
*/
func (l *ThresholdPrivKeyGenerator) GeneratePhase1PubKeyProof() (msg *models.KeyGenBroadcastMessage1, err error) {
	ui, dk := createKeys()
	p := &models.PrivateKeyInfo{
		Key:             l.PrivateKeyID,
		UI:              ui,
		PaillierPrivkey: dk,
		Status:          models.PrivateKeyNegotiateStatusInit,
	}
	err = l.db.NewPrivateKeyInfo(p)
	if err != nil {
		return
	}
	//用dlogproof 传播我自己的公钥片,包含相关的零知识证明
	proof := proofs.Prove(ui)
	msg = &models.KeyGenBroadcastMessage1{Proof: proof}
	return
}

/*
ReceivePhase1PubKeyProof phase 1.2 接受其他公证人传递过来的公钥片证明信息,
如果凑齐了所有公证人(params.ShareCount)的公钥片信息,那么就可以组合出来此次协商最终的公钥.
但是到目前为止没有一个公证人知道最终公钥对应的私钥片
*/
func (l *ThresholdPrivKeyGenerator) ReceivePhase1PubKeyProof(m *models.KeyGenBroadcastMessage1, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
	if p.PubKeysProof1 != nil {
		_, ok := p.PubKeysProof1[index]
		if ok {
			err = fmt.Errorf("pubkey roof for %d already exist", index)
			return
		}
	}
	if !proofs.Verify(m.Proof) {
		err = fmt.Errorf("pubkey roof for %d verify not pass", index)
		return
	}
	if p.PubKeysProof1 == nil {
		p.PubKeysProof1 = make(map[int]*models.KeyGenBroadcastMessage1)
	}
	p.PubKeysProof1[index] = m
	err = l.db.KGUpdatePubKeysProof1(p)
	if err != nil {
		return
	}
	//除了自己以外的因素都凑齐了
	if len(p.PubKeysProof1) == params.ShareCount-1 {
		x, y := share.S.ScalarBaseMult(p.UI.Bytes())
		for _, m := range p.PubKeysProof1 {
			x, y = share.PointAdd(x, y, m.Proof.PK.X, m.Proof.PK.Y)
		}
		//所有人公证人都掌握的总公钥
		p.PublicKeyX = x
		p.PublicKeyY = y
		err = l.db.KGUpdateTotalPubKey(p)
	}
	finish = len(p.PubKeysProof1) == params.ShareCount-1
	return
}

func createCommitmentWithUserDefinedRandomNess(message *big.Int, blindingFactor *big.Int) *big.Int {
	hash := utils.Sha256(message.Bytes(), blindingFactor.Bytes())
	b := new(big.Int)
	b.SetBytes(hash[:])
	return b
}

/*
GeneratePhase2PaillierKeyProof phase 2.1 广播给其他所有公证人的同态加密公钥以及Proof
*/
func (l *ThresholdPrivKeyGenerator) GeneratePhase2PaillierKeyProof() (msg *models.KeyGenBroadcastMessage2, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
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
	p.PaillierKeysProof2 = make(map[int]*models.KeyGenBroadcastMessage2)
	p.PaillierKeysProof2[l.selfNotaryID] = msg
	//保存自己的信息到数据库中,方便后续计算使用
	err = l.db.KGUpdatePaillierKeysProof2(p)
	return
}

//ReceivePhase2PaillierPubKeyProof phase 2.2 收到其他公证人的同态加密公钥信息
func (l *ThresholdPrivKeyGenerator) ReceivePhase2PaillierPubKeyProof(m *models.KeyGenBroadcastMessage2,
	index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
	//不接受重复的消息
	_, ok := p.PaillierKeysProof2[index]
	if ok {
		err = fmt.Errorf("paillier pubkey roof for %d already exist", index)
		return
	}
	//对方证明私钥我确实知道,虽然没有告诉我
	if !m.CorrectKeyProof.Verify(m.PaillierPubkey) {
		err = fmt.Errorf("paillier pubkey roof for %d verify not pass", index)
		return
	}
	if createCommitmentWithUserDefinedRandomNess(p.PubKeysProof1[index].Proof.PK.X, m.BlindFactor).Cmp(m.Com) != 0 {
		err = fmt.Errorf("blind factor error for %d", index)
		return
	}
	p.PaillierKeysProof2[index] = m
	err = l.db.KGUpdatePaillierKeysProof2(p)
	if err != nil {
		return
	}
	//所有因素都凑齐了
	finish = len(p.PaillierKeysProof2) == params.ShareCount
	return
}

/*
GeneratePhase3SecretShare phase 3.1基于第一步的随机数生成SecretShares,
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
func (l *ThresholdPrivKeyGenerator) GeneratePhase3SecretShare() (msgs map[int]*models.KeyGenBroadcastMessage3, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}

	vss, secretShares := feldman.Share(params.ThresholdCount, params.ShareCount, p.UI)
	msg := &models.KeyGenBroadcastMessage3{
		Vss:         vss,
		SecretShare: secretShares[l.selfNotaryID],
		Index:       l.selfNotaryID,
	}
	p.SecretShareMessage3 = make(map[int]*models.KeyGenBroadcastMessage3)
	p.SecretShareMessage3[l.selfNotaryID] = msg
	err = l.db.KGUpdateSecretShareMessage3(p)
	if err != nil {
		return
	}
	msgs = make(map[int]*models.KeyGenBroadcastMessage3)
	for i := 0; i < params.ShareCount; i++ {
		if i == l.selfNotaryID {
			continue
		}
		msg := &models.KeyGenBroadcastMessage3{
			Vss:         vss,
			SecretShare: secretShares[i],
			Index:       i,
		}
		msgs[i] = msg
	}
	return
}

//ReceivePhase3SecretShare 接收来自其他公证人定向分发给我的secret share.
func (l *ThresholdPrivKeyGenerator) ReceivePhase3SecretShare(msg *models.KeyGenBroadcastMessage3, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
	_, ok := p.SecretShareMessage3[index]
	if ok {
		err = fmt.Errorf("secret shares for %d already received", index)
		return
	}
	//vss的第0号证据对应的就是第一步随机数对应的公钥片,必须用的是那一个
	if !EqualGE(p.PubKeysProof1[index].Proof.PK, msg.Vss.Commitments[0]) {
		err = fmt.Errorf("pubkey not match for %d", index)
		return
	}
	//必须经过feldman vss验证,符合规则.如果某个公证人给的secret share发错了,比如s11发送给了2号公证人,这里会检测出错误.
	if !msg.Vss.ValidateShare(msg.SecretShare, l.selfNotaryID+1) {
		err = fmt.Errorf("secret share error for %d", index)
		return
	}
	p.SecretShareMessage3[index] = msg
	err = l.db.KGUpdateSecretShareMessage3(p)
	if err != nil {
		return
	}
	finish = len(p.SecretShareMessage3) == params.ShareCount
	//所有因素都凑齐了.可以计算我自己持有的私钥片 sum(f(i)) f(i)解释:公证人生成了n个多项式 ,把所有第i个多项式的结果相加)
	if len(p.SecretShareMessage3) == params.ShareCount {
		p.XI = share.SPrivKey{D: new(big.Int)}
		for _, s := range p.SecretShareMessage3 {
			share.ModAdd(p.XI, s.SecretShare)
		}
	}
	err = l.db.KGUpdateXI(p)
	return
}

/*
GeneratePhase4PubKeyProof phase4.1 生成我自己对我持有私钥片的证明,
todo 目前问题就在于对方无法验证对方给我的XI是真实有效的,和公钥是关联在一起的真实私钥片,
有可能对方在GeneratePhase3SecretShare 就广播一个错误的secret share,但是我也无法证明
*/
func (l *ThresholdPrivKeyGenerator) GeneratePhase4PubKeyProof() (msg *models.KeyGenBroadcastMessage4, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
	if p.XI.D == nil {
		panic("must wait for phase 3 finish")
	}
	proof := proofs.Prove(p.XI)
	msg = &models.KeyGenBroadcastMessage4{Proof: proof}
	p.LastPubkeyProof4 = make(map[int]*models.KeyGenBroadcastMessage4)
	p.LastPubkeyProof4[l.selfNotaryID] = msg
	err = l.db.KGUpdateLastPubKeyProof4(p)
	return
}

/*
ReceivePhase4VerifyTotalPubKey phase4.2 接收对方持有的Xi的证明,并验证其有效性. 但是我怎么知道他是有效的呢? 目前为止是没有办法的.
*/
func (l *ThresholdPrivKeyGenerator) ReceivePhase4VerifyTotalPubKey(msg *models.KeyGenBroadcastMessage4, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKeyID)
	if err != nil {
		return
	}
	_, ok := p.LastPubkeyProof4[index]
	if ok {
		err = fmt.Errorf("last pubkey for %d already exist", index)
		return
	}
	if !proofs.Verify(msg.Proof) {
		err = fmt.Errorf("last pubkey verify error for %d", index)
		return
	}
	p.LastPubkeyProof4[index] = msg
	err = l.db.KGUpdateLastPubKeyProof4(p)
	if err != nil {
		return
	}

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
	}

	finish = len(p.LastPubkeyProof4) == params.ShareCount
	p.Status = models.PrivateKeyNegotiateStatusFinished
	err = l.db.KGUpdateKeyGenStatus(p)
	return
}

// EqualGE :
func EqualGE(pubGB *share.SPubKey, mtaGB *share.SPubKey) bool {
	return pubGB.X.Cmp(mtaGB.X) == 0 && pubGB.Y.Cmp(mtaGB.Y) == 0
}
