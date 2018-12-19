package mecdsa

import (
	"errors"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/params"

	"math/big"

	"bytes"

	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/goutils"
	"github.com/nkbai/log"
)

/*
DistributedSignMessage 由某个公证人主导发起,其他一组公证人参与的,相互协调最终生成有效签名的过程.
*/
type DistributedSignMessage struct {
	db                 *models.DB
	srv                *NotaryService
	Key                common.Hash               //此次签名的唯一key,由签名主导公证人指定
	PrivateKey         common.Hash               //此次签名用到的分布式私钥, 在数据库中的key
	Message            []byte                    //此次签名的消息
	S                  []int                     //由签名主导公证人指定的此次签名参与的公证人.
	XI                 share.SPrivKey            //协商好的私钥片
	PaillierPubKeys    map[int]*proofs.PublicKey //其他公证人的同态加密公钥
	PaillierPrivateKey *proofs.PrivateKey        //我的同态加密私钥
	PublicKey          *share.SPubKey            //上次协商生成的总公钥
	Vss                *feldman.VerifiableSS     //上次协商生成的feldman vss
	L                  *models.SignMessage       //此次签名生成过程中需要保存到数据库的信息
	index              int
}

/*
NewDistributedSignMessage 一开始就要确定哪些公证人参与此次签名生成,
人数t > ThresholdCount && t <= ShareCount
指出要签名的交易,公证人应该对此交易做校验,是否是一个合法的交易
*/
func NewDistributedSignMessage(db *models.DB, srv *NotaryService, message []byte, key common.Hash, privateKey common.Hash, s []int) (l *DistributedSignMessage, err error) {
	if len(s) <= params.ThresholdCount {
		err = fmt.Errorf("candidates notary too less")
		return
	}
	l = &DistributedSignMessage{
		db:         db,
		srv:        srv,
		Key:        key,
		Message:    message,
		PrivateKey: privateKey,
		S:          s,
		index:      srv.NotaryShareArg.Index,
	}
	l2 := &models.SignMessage{
		Key:            l.Key,
		UsedPrivateKey: l.PrivateKey,
		Message:        l.Message,
		S:              s,
	}
	err = l.db.NewLockedout(l2)
	if err != nil {
		return nil, err
	}
	err = l.loadLockout()
	if err != nil {
		return nil, err
	}
	l.createSignKeys()
	return
}

// 从数据库中载入相关信息,可以做缓存
func (l *DistributedSignMessage) loadLockout() error {
	//if l.L != nil {
	//	return nil
	//}
	p, err := l.db.LoadPrivatedKeyInfo(l.PrivateKey)
	if err != nil {
		return err
	}
	l.XI = p.XI
	l.PaillierPrivateKey = p.PaillierPrivkey
	l.PaillierPubKeys = make(map[int]*proofs.PublicKey)
	for k, v := range p.PaillierKeysProof2 {
		l.PaillierPubKeys[k] = v.PaillierPubkey
	}
	l.PublicKey = &share.SPubKey{p.PublicKeyX, p.PublicKeyY}
	l.Vss = p.SecretShareMessage3[l.srv.NotaryShareArg.Index].Vss

	l.L, err = l.db.LoadLockout(l.Key)
	if err != nil {
		return err
	}
	return nil
}

/*
		//参数：自己的私钥片、系数点乘集合和我的多项式y集合、签名人的原始编号、所有签名人的编号
//==>每个公证人的公私钥、公钥片、{{{t,n},t+1个系数点乘G的结果(c1...c2)},y1...yn}
*/
func (l *DistributedSignMessage) createSignKeys() {
	lambdaI := l.Vss.MapShareToNewParams(l.srv.NotaryShareArg.Index, l.S) //lamda_i 解释：通过lamda_i对原所有证明人的群来映射出签名者群
	wi := share.ModMul(lambdaI, l.XI)                                     //wi： 我原来的编号在对应签名群编号的映射关系 ，原来我是xi(私钥片) 现在是wi（我在新的签名群中的私钥片）
	gwiX, gwiY := share.S.ScalarBaseMult(wi.Bytes())                      //我在签名群中的公钥片
	gammaI := share.RandomPrivateKey()                                    //临时私钥
	gGammaIX, gGammaIY := share.S.ScalarBaseMult(gammaI.Bytes())          //临时公钥
	l.L.SignedKey = &models.SignedKey{
		WI:      wi,
		Gwi:     &share.SPubKey{gwiX, gwiY},
		KI:      share.RandomPrivateKey(),
		GammaI:  gammaI,
		GGammaI: &share.SPubKey{gGammaIX, gGammaIY},
	}
}

/*
GeneratePhase1Broadcast 确定此次签名所用临时私钥,不能再换了.
*/
func (l *DistributedSignMessage) GeneratePhase1Broadcast() (msg *models.SignBroadcastPhase1, err error) {
	blindFactor := share.RandomBigInt()
	gGammaIX, _ := share.S.ScalarBaseMult(l.L.SignedKey.GammaI.Bytes())
	com := createCommitmentWithUserDefinedRandomNess(gGammaIX, blindFactor)
	msg = &models.SignBroadcastPhase1{
		Com:         com,
		BlindFactor: blindFactor,
	}
	l.L.Phase1BroadCast = make(map[int]*models.SignBroadcastPhase1)
	l.L.Phase1BroadCast[l.index] = msg
	err = l.db.UpdateLockout(l.L)
	return
}

//ReceivePhase1Broadcast 收集此次签名中,其他公证人所用临时公钥,保证在后续步骤中不会被替换
func (l *DistributedSignMessage) ReceivePhase1Broadcast(msg *models.SignBroadcastPhase1, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	_, ok := l.L.Phase1BroadCast[index]
	if ok {
		err = fmt.Errorf("phase1 broadcast for %d already exist", index)
		return
	}
	l.L.Phase1BroadCast[index] = msg
	err = l.db.UpdateLockout(l.L)
	finish = len(l.L.Phase1BroadCast) == len(l.S)
	return
}

//p2p告知E(ki) 签名发起人，告诉其他所有人 步骤1 生成ki， 保证后续ki不会发生变化
func newMessageA(ki share.SPrivKey, paillierPubKey *proofs.PublicKey) (*models.MessageA, error) {
	ca, err := proofs.Encrypt(paillierPubKey, ki.Bytes())
	if err != nil {
		return nil, err
	}
	return &models.MessageA{ca}, nil
}

/*
GeneratePhase2MessageA p2p告知E(ki) 签名发起人，告诉其他所有人 步骤1 生成ki， 保证后续ki不会发生变化,
同时其他人是无法获取到ki,只有自己知道自己的ki
*/
func (l *DistributedSignMessage) GeneratePhase2MessageA() (msg *models.MessageA, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	ma, err := newMessageA(l.L.SignedKey.KI, &l.PaillierPrivateKey.PublicKey)
	if err != nil {
		return
	}
	l.L.MessageA = ma
	l.L.Phase2MessageB = make(map[int]*models.MessageBPhase2)
	l.L.AlphaGamma = make(map[int]share.SPrivKey)
	l.L.AlphaWI = make(map[int]share.SPrivKey)
	err = l.db.UpdateLockout(l.L)
	return ma, err
}

/*NewMessageB
//我(bob)签名key中的临时私钥gammaI、(alice)paillier公钥、(alice)paillier公钥加密ki的结果
//返回：cB和两个证明(其他人验证)
我（alice）收到其他人（bob)给我的ma以后， 立即计算mb和两个证明，发送给bob
gammaI: alice的临时私钥
paillierPubKey: bob的同态加密公钥
ca: bob发送给alice的E(Ki)
*/
func NewMessageB(gammaI share.SPrivKey, paillierPubKey *proofs.PublicKey, ca *models.MessageA) (*models.MessageB, error) {
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
func (l *DistributedSignMessage) ReceivePhase2MessageA(msg *models.MessageA, index int) (mb *models.MessageBPhase2, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	mgGamma, err := NewMessageB(l.L.SignedKey.GammaI, l.PaillierPubKeys[index], msg)
	if err != nil {
		return
	}
	mbw, err := NewMessageB(l.L.SignedKey.WI, l.PaillierPubKeys[index], msg)
	if err != nil {
		return
	}
	mb = &models.MessageBPhase2{
		MessageBGamma: mgGamma,
		MessageBWi:    mbw,
	}
	//l.L.Phase2MessageB=make(map[int]*models.MessageBPhase2)
	///l.L.Phase2MessageB[l.index]=
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

//ReceivePhase2MessageB 收集并验证MessageB证明信息
func (l *DistributedSignMessage) ReceivePhase2MessageB(msg *models.MessageBPhase2, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	alphaijGamma, err := verifyProofsGetAlpha(msg.MessageBGamma, l.PaillierPrivateKey, l.L.SignedKey.KI)
	if err != nil {
		return
	}
	alphaijWi, err := verifyProofsGetAlpha(msg.MessageBWi, l.PaillierPrivateKey, l.L.SignedKey.KI)
	if err != nil {
		return
	}
	//if !EqualGE(msg.MessageBWi.BProof.PK, l.L.SignedKey.Gwi) { 这里应该等于index的gwi,而不是l.index的gwi
	//	panic("not equal")
	//}
	l.L.Phase2MessageB[index] = msg
	l.L.AlphaGamma[index] = alphaijGamma
	l.L.AlphaWI[index] = alphaijWi
	err = l.db.UpdateLockout(l.L)
	if err != nil {
		return
	}
	finish = len(l.L.Phase2MessageB) == len(l.S)-1
	return
}

func (l *DistributedSignMessage) phase2DeltaI() share.SPrivKey {
	var k *models.SignedKey
	k = l.L.SignedKey
	if len(l.L.AlphaGamma) != len(l.S)-1 {
		panic("arg error")
	}
	//kiGammaI=ki * gammI+Sum(alpha_vec) +Sum(beta_vec)
	kiGammaI := k.KI.Clone()
	share.ModMul(kiGammaI, k.GammaI)
	for _, i := range l.S {
		if i == l.index {
			continue
		}
		share.ModAdd(kiGammaI, l.L.AlphaGamma[i])
		share.ModAdd(kiGammaI, l.L.Phase2MessageB[i].MessageBGamma.Beta)
	}
	return kiGammaI
}
func (l *DistributedSignMessage) phase2SigmaI() share.SPrivKey {
	if len(l.L.AlphaWI) != len(l.S)-1 {
		panic("length error")
	}
	kiwi := l.L.SignedKey.KI.Clone()
	share.ModMul(kiwi, l.L.SignedKey.WI)
	//todo vij=vji ?
	for _, i := range l.S {
		if i == l.index {
			continue
		}
		share.ModAdd(kiwi, l.L.AlphaWI[i])
		share.ModAdd(kiwi, l.L.Phase2MessageB[i].MessageBWi.Beta)
	}
	return kiwi
}

/*
GeneratePhase3DeltaI 依据上一步协商信息,生成我自己的DeltaI,然后广播给所有其他人,需要这些参与公证人得到完整的的Delta
但是生成的SigmaI自己保留,在生成自己的签名片的时候使用
*/
func (l *DistributedSignMessage) GeneratePhase3DeltaI() (msg *models.DeltaPhase3, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	deltaI := l.phase2DeltaI()
	sigmaI := l.phase2SigmaI()
	l.L.Sigma = sigmaI
	l.L.Delta = make(map[int]share.SPrivKey)
	l.L.Delta[l.index] = deltaI
	err = l.db.UpdateLockout(l.L)
	if err != nil {
		return
	}
	msg = &models.DeltaPhase3{deltaI}
	return
}

//ReceivePhase3DeltaI 收集所有的deltaI
func (l *DistributedSignMessage) ReceivePhase3DeltaI(msg *models.DeltaPhase3, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	_, ok := l.L.Delta[index]
	if ok {
		err = fmt.Errorf("ReceivePhase3DeltaI for %d already exist", index)
		return
	}
	l.L.Delta[index] = msg.Delta
	err = l.db.UpdateLockout(l.L)
	finish = len(l.L.Delta) == len(l.S)
	return
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
	return &share.SPubKey{rx, ry}, nil
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

//GeneratePhase4R 所有公证人都应该得到相同的R,其中R.X就是最后签名(r,s,v)中的r
func (l *DistributedSignMessage) GeneratePhase4R() (R *share.SPubKey, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	delta := phase3ReconstructDelta(l.L.Delta)
	//少一个自己的MessageB todo fixme 这里面需要有更好的方式来生成数据,后续必须优化
	mgGamma, err := NewMessageB(l.L.SignedKey.GammaI, &l.PaillierPrivateKey.PublicKey, l.L.MessageA)
	if err != nil {
		return
	}
	l.L.Phase2MessageB[l.index] = &models.MessageBPhase2{mgGamma, nil}
	R, err = phase4(delta, l.L.Phase2MessageB, l.L.Phase1BroadCast)
	if err != nil {
		return
	}
	l.L.R = R
	err = l.db.UpdateLockout(l.L)
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
func phase5aBroadcast5bZkproof(l *models.LocalSignature) (*models.Phase5Com1, *models.Phase5ADecom1, *proofs.HomoELGamalProof) {
	blindFactor := share.RandomBigInt()
	//Ai=g^^rho_i
	aix, aiy := share.S.ScalarBaseMult(l.RhoI.Bytes())
	lIRhoI := l.LI.Clone()
	share.ModMul(lIRhoI, l.RhoI)
	//Bi=G*lIRhoI
	bix, biy := share.S.ScalarBaseMult(lIRhoI.Bytes())
	//vi=R*si+G*li
	tx, ty := share.S.ScalarMult(l.R.X, l.R.Y, l.SI.Bytes()) //R^^si
	vix, viy := share.S.ScalarBaseMult(l.LI.Bytes())         //g^^li
	vix, viy = share.PointAdd(vix, viy, tx, ty)

	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{
		{vix, viy}, {aix, aiy}, {bix, biy},
	})
	com := createCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)

	//proof是5b的zkp构造
	witness := proofs.NewHomoElGamalWitness(l.LI, l.SI) //li si
	delta := &proofs.HomoElGamalStatement{
		G: share.NewGE(aix, aiy),               //Ai
		H: share.NewGE(l.R.X, l.R.Y),           //R
		Y: share.NewGE(share.S.Gx, share.S.Gy), //g
		D: share.NewGE(vix, viy),               //Vi
		E: share.NewGE(bix, biy),               //Bi
	}
	//证明提供的是正确的si???
	proof := proofs.CreateHomoELGamalProof(witness, delta)
	return &models.Phase5Com1{com},
		&models.Phase5ADecom1{
			Vi:          share.NewGE(vix, viy),
			Ai:          share.NewGE(aix, aiy),
			Bi:          share.NewGE(bix, biy),
			BlindFactor: blindFactor,
		},
		proof
}

/*
GeneratePhase5a5bZkProof 从此步骤开始,互相不断交换信息来判断对方能够生成正确的Si(也就是签名片),
如果所有参与者都能生成最终的签名片,那么我才能把自己的签名片告诉对方.
si的累加和就是签名(r,s,v)中的s
*/
func (l *DistributedSignMessage) GeneratePhase5a5bZkProof() (msg *models.Phase5A, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	messageHash := utils.Sha256(l.Message)
	messageBN := new(big.Int).SetBytes(messageHash[:])
	localSignature := phase5LocalSignature(l.L.SignedKey.KI, messageBN, l.L.R, l.L.Sigma, l.PublicKey)
	phase5Com, phase5ADecom, helgamalProof := phase5aBroadcast5bZkproof(localSignature)
	msg = &models.Phase5A{phase5Com, phase5ADecom, helgamalProof}
	l.L.Phase5A = make(map[int]*models.Phase5A)
	l.L.Phase5A[l.index] = msg
	l.L.LocalSignature = localSignature
	err = l.db.UpdateLockout(l.L)
	return
}

func (l *DistributedSignMessage) ReceivePhase5A5BProof(msg *models.Phase5A, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	_, ok := l.L.Phase5A[index]
	if ok {
		err = fmt.Errorf("ReceivePhase5A5BProof already exist for %d", index)
		return
	}
	delta := &proofs.HomoElGamalStatement{
		G: msg.Phase5ADecom1.Ai,
		H: l.L.R,
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
		err = errors.New("invalid com")
		return
	}
	l.L.Phase5A[index] = msg
	finish = len(l.L.Phase5A) == len(l.S)
	err = l.db.UpdateLockout(l.L)
	return
}

func phase5c(l *models.LocalSignature, deCommitments []*models.Phase5ADecom1,
	vi *share.SPubKey) (*models.Phase5Com2, *models.Phase5DDecom2, error) {

	//从广播的commit(Ci,Di)得到vi,ai
	v := vi.Clone()
	for i := 0; i < len(deCommitments); i++ {
		v.X, v.Y = share.PointAdd(v.X, v.Y, deCommitments[i].Vi.X, deCommitments[i].Vi.Y)
	}
	a := deCommitments[0].Ai.Clone()
	for i := 1; i < len(deCommitments); i++ {
		a.X, a.Y = share.PointAdd(a.X, a.Y, deCommitments[i].Ai.X, deCommitments[i].Ai.Y)
	}
	r := share.BigInt2PrivateKey(l.R.X)
	yrx, yry := share.S.ScalarMult(l.Y.X, l.Y.Y, r.Bytes())
	m := share.BigInt2PrivateKey(l.M)
	//Vi之积×g^(-m)*y^(-r)
	gmx, gmy := share.S.ScalarBaseMult(m.Bytes())
	v.X, v.Y = share.PointSub(v.X, v.Y, gmx, gmy)
	v.X, v.Y = share.PointSub(v.X, v.Y, yrx, yry)
	//Ui=V * rhoi
	uix, uiy := share.S.ScalarMult(v.X, v.Y, l.RhoI.Bytes())
	//Ti=A * li
	tix, tiy := share.S.ScalarMult(a.X, a.Y, l.LI.Bytes())

	//commit(Ui ,Ti)，广播出去
	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{
		{uix, uiy},
		{tix, tiy},
	})
	blindFactor := share.RandomBigInt()
	com := createCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)
	return &models.Phase5Com2{com},
		&models.Phase5DDecom2{
			Ui:          &share.SPubKey{uix, uiy}, //Ci
			Ti:          &share.SPubKey{tix, tiy}, //Di
			BlindFactor: blindFactor,
		},
		nil
}

//GeneratePhase5CProof  fixme 提供一个好的注释
func (l *DistributedSignMessage) GeneratePhase5CProof() (msg *models.Phase5C, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	if len(l.L.Phase5A) != len(l.S) {
		panic("cannot genrate5c until all 5b proof received")
	}
	var decomVec []*models.Phase5ADecom1
	for i, m := range l.L.Phase5A {
		if i == l.index { //phase5c 不应该包括自己的decommitment
			continue
		}
		decomVec = append(decomVec, m.Phase5ADecom1)
	}
	phase5com2, phase5decom2, err := phase5c(l.L.LocalSignature, decomVec, l.L.Phase5A[l.index].Phase5ADecom1.Vi)
	if err != nil {
		return
	}
	msg = &models.Phase5C{phase5com2, phase5decom2}
	l.L.Phase5C = make(map[int]*models.Phase5C)
	l.L.Phase5C[l.index] = msg
	err = l.db.UpdateLockout(l.L)
	return
}

//ReceivePhase5cProof   fixme 暂时没有好的解释
func (l *DistributedSignMessage) ReceivePhase5cProof(msg *models.Phase5C, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	//验证5c的hash(ui,ti)=5c的Ci

	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{msg.Phase5DDecom2.Ui, msg.Phase5DDecom2.Ti})
	inputhash.D = createCommitmentWithUserDefinedRandomNess(inputhash.D, msg.Phase5DDecom2.BlindFactor)
	if inputhash.D.Cmp(msg.Phase5Com2.Com) != 0 {
		err = errors.New("invalid com")
		return
	}
	_, ok := l.L.Phase5C[index]
	if ok {
		err = fmt.Errorf("ReceivePhase5cProof for %d already exist", index)
		return
	}
	l.L.Phase5C[index] = msg
	err = l.db.UpdateLockout(l.L)
	finish = len(l.L.Phase5C) == len(l.S)
	return
}

/*
Generate5dProof
接受所有签名人的si的广播，有可能某个公证人会保留信息，最终生成有效的签名，私自保留下来,但是不告诉其他人自己的si是多少.
但是这种情况其他公证人可以知道,没有收到某个公证人的si
*/
func (l *DistributedSignMessage) Generate5dProof() (si share.SPrivKey, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	var deCommitments2 []*models.Phase5DDecom2
	var commitments2 []*models.Phase5Com2
	var deCommitments1 []*models.Phase5ADecom1
	for i, m := range l.L.Phase5C {
		deCommitments2 = append(deCommitments2, m.Phase5DDecom2)
		commitments2 = append(commitments2, m.Phase5Com2)
		deCommitments1 = append(deCommitments1, l.L.Phase5A[i].Phase5ADecom1)
	}
	//校验收到的关于si的信息,都是真实的,正确的,只有通过了,才能把自己的si告诉给其他公证人
	si, err = phase5d(l.L.LocalSignature, deCommitments2, commitments2, deCommitments1)
	if err != nil {
		return
	}
	l.L.Phase5D = make(map[int]share.SPrivKey)
	l.L.Phase5D[l.index] = si
	err = l.db.UpdateLockout(l.L)
	return
}

/*
 校验收到的关于si的信息,都是真实的,正确的,只有通过了,才能把自己的si告诉给其他公证人
*/
func phase5d(l *models.LocalSignature, deCommitments2 []*models.Phase5DDecom2,
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
	//用于比较 Ui 和 (5c的ti + 5a的bi)连加 是否相等
	for i := 0; i < len(commitments2); i++ {
		biasedSumTbX, biasedSumTbY = share.PointSub(
			biasedSumTbX, biasedSumTbY,
			deCommitments2[i].Ui.X, deCommitments2[i].Ui.Y,
		)
	}
	log.Trace(fmt.Sprintf("(gx,gy)=(%s,%s)", share.S.Gx.Text(16), share.S.Gy.Text(16)))
	log.Trace(fmt.Sprintf("(tbx,tby)=(%s,%s)", biasedSumTbX.Text(16), biasedSumTbY.Text(16)))
	if share.S.Gx.Cmp(biasedSumTbX) == 0 &&
		share.S.Gy.Cmp(biasedSumTbY) == 0 {
		return l.SI.Clone(), nil
	}
	return share.PrivKeyZero, errors.New("invalid key")
}

//RecevieSI 收集签名片,收集齐以后就可以得到完整的签名.所有公证人都应该得到有效的签名.
func (l *DistributedSignMessage) RecevieSI(si share.SPrivKey, index int) (signature []byte, finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	if _, ok := l.L.Phase5D[index]; ok {
		err = fmt.Errorf("si for %d already received", index)
		return
	}
	l.L.Phase5D[index] = si
	err = l.db.UpdateLockout(l.L)
	if err != nil {
		return
	}
	finish = len(l.L.Phase5D) == len(l.L.S)
	if !finish {
		return
	}
	s := l.L.LocalSignature.SI.Clone()
	//所有人的的si，包括自己
	for i, si := range l.L.Phase5D {
		if i == l.index {
			continue
		}
		share.ModAdd(s, si)
	}
	r := share.BigInt2PrivateKey(l.L.R.X)
	signature, verifyResult := verify(s, r, l.L.LocalSignature.Y, l.L.LocalSignature.M)
	if !verifyResult {
		err = errors.New("invilad signature")
	}
	return
}

/*
由于目前的信息只有r,s,没有v,应该有办法直接得到v是0还是1,
但是目前只能通过尝试.
*/
func verify(s, r share.SPrivKey, y *share.SPubKey, message *big.Int) ([]byte, bool) {
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
	//缺少信息v(0或者1),所以尝试这两种情况,只要有一个可以被矿工验证,那就认为得到了有效签名.
	key, _ := crypto.GenerateKey()
	pubkey := key.PublicKey
	pubkey.X = y.X
	pubkey.Y = y.Y
	addr := crypto.PubkeyToAddress(pubkey)
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(r.D))
	buf.Write(utils.BigIntTo32Bytes(s.D))
	buf.Write([]byte{0})
	bs := buf.Bytes()
	h := common.Hash{}
	h.SetBytes(message.Bytes())
	pubkeybin, err := crypto.Ecrecover(h[:], bs)
	if err != nil {
		return nil, false
	}
	pubkey2, _ := crypto.UnmarshalPubkey(pubkeybin)
	addr2 := crypto.PubkeyToAddress(*pubkey2)
	if addr2 == addr {
		return bs, true
	}
	bs[64] = 1
	pubkeybin, err = crypto.Ecrecover(h[:], bs)
	if err != nil {
		return nil, false
	}
	pubkey2, _ = crypto.UnmarshalPubkey(pubkeybin)
	addr2 = crypto.PubkeyToAddress(*pubkey2)
	if addr2 == addr {
		return bs, true
	}
	return nil, false

}
