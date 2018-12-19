package mecdsa

import (
	"errors"
	"fmt"

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

type Lockout struct {
	db                 *models.DB
	srv                *NotaryService
	Key                common.Hash
	PrivateKey         common.Hash
	Message            []byte
	S                  []int
	XI                 share.SPrivKey //协商好的私钥片
	PaillierPubKeys    map[int]*proofs.PublicKey
	PaillierPrivateKey *proofs.PrivateKey
	PublicKey          *share.SPubKey
	Vss                *feldman.VerifiableSS
	L                  *models.Lockout
	index              int
}

func NewLockout(db *models.DB, srv *NotaryService, message []byte, key common.Hash, privateKey common.Hash, s []int) (l *Lockout, err error) {

	l = &Lockout{
		db:         db,
		srv:        srv,
		Key:        key,
		Message:    message,
		PrivateKey: privateKey,
		S:          s,
		index:      srv.NotaryShareArg.Index,
	}
	l2 := &models.Lockout{
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

func (l *Lockout) loadLockout() error {
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
func (l *Lockout) createSignKeys() {
	li := l.Vss.MapShareToNewParams(l.srv.NotaryShareArg.Index, l.S) //lamda_i 解释：通过lamda_i对原所有证明人的群来映射出签名者群
	wi := share.ModMul(li, l.XI)                                     //wi： 我原来的编号在对应签名群编号的映射关系 ，原来我是xi(私钥片) 现在是wi（我在新的签名群中的私钥片）
	gwiX, gwiY := share.S.ScalarBaseMult(wi.Bytes())                 //我在签名群中的公钥片
	gammaI := share.RandomPrivateKey()                               //临时私钥
	gGammaIX, gGammaIY := share.S.ScalarBaseMult(gammaI.Bytes())     //临时公钥
	l.L.SignedKey = &models.SignedKey{
		WI:      wi,
		Gwi:     &share.SPubKey{gwiX, gwiY},
		KI:      share.RandomPrivateKey(),
		GammaI:  gammaI,
		GGammaI: &share.SPubKey{gGammaIX, gGammaIY},
	}
}

func (l *Lockout) GeneratePhase1Broadcast() (msg *models.SignBroadcastPhase1, err error) {
	blindFactor := share.RandomBigInt()
	gGammaIX, _ := share.S.ScalarBaseMult(l.L.SignedKey.GammaI.Bytes())
	com := CreateCommitmentWithUserDefinedRandomNess(gGammaIX, blindFactor)
	msg = &models.SignBroadcastPhase1{
		Com:         com,
		BlindFactor: blindFactor,
	}
	l.L.Phase1BroadCast = make(map[int]*models.SignBroadcastPhase1)
	l.L.Phase1BroadCast[l.index] = msg
	err = l.db.UpdateLockout(l.L)
	return
}

func (l *Lockout) ReceivePhase1Broadcast(msg *models.SignBroadcastPhase1, index int) (finish bool, err error) {
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

//const
func newMessageA(ecdsaPrivateKey share.SPrivKey, paillierPubKey *proofs.PublicKey) (*models.MessageA, error) {
	ca, err := proofs.Encrypt(paillierPubKey, ecdsaPrivateKey.Bytes())
	if err != nil {
		return nil, err
	}
	return &models.MessageA{ca}, nil
}
func (l *Lockout) GeneratePhase2MessageA() (msg *models.MessageA, err error) {
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

//NewMessageB
//我(bob)签名key中的临时私钥gammaI、(alice)paillier公钥、(alice)paillier公钥加密ki的结果
//返回：cB和两个证明(其他人验证)
func NewMessageB(gammaI share.SPrivKey, paillierPubKey *proofs.PublicKey, ca *models.MessageA) (*models.MessageB, error) {
	betaTag := share.RandomBigInt()
	//todo fixme bai
	//betaTag = big.NewInt(39)
	betaTagPrivateKey := share.BigInt2PrivateKey(betaTag.Mod(betaTag, share.S.N))
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
func (l *Lockout) ReceivePhase2MessageA(msg *models.MessageA, index int) (mb *models.MessageBPhase2, err error) {
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

//const
func VerifyProofsGetAlpha(m *models.MessageB, dk *proofs.PrivateKey, a share.SPrivKey) (share.SPrivKey, error) {
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
func (l *Lockout) ReceivePhase2MessageB(msg *models.MessageBPhase2, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	alphaijGamma, err := VerifyProofsGetAlpha(msg.MessageBGamma, l.PaillierPrivateKey, l.L.SignedKey.KI)
	if err != nil {
		return
	}
	alphaijWi, err := VerifyProofsGetAlpha(msg.MessageBWi, l.PaillierPrivateKey, l.L.SignedKey.KI)
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

//const
func (l *Lockout) phase2DeltaI() share.SPrivKey {
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
func (l *Lockout) phase2SigmaI() share.SPrivKey {
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
func (l *Lockout) GeneratePhase3DeltaI() (msg *models.DeltaPhase3, err error) {
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

func (l *Lockout) ReceivePhase3DeltaI(msg *models.DeltaPhase3, index int) (finish bool, err error) {
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

//const
func phase4(delta share.SPrivKey,
	gammaIProof map[int]*models.MessageBPhase2,
	phase1Broadcast map[int]*models.SignBroadcastPhase1) (*share.SPubKey, error) {
	if len(gammaIProof) != len(phase1Broadcast) {
		panic("length must equal")
	}
	for i, p := range gammaIProof {
		if CreateCommitmentWithUserDefinedRandomNess(p.MessageBGamma.BProof.PK.X, phase1Broadcast[i].BlindFactor).Cmp(phase1Broadcast[i].Com) == 0 {
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

func (l *Lockout) GeneratePhase4R() (R *share.SPubKey, err error) {
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

//const
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
	com := CreateCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)

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
func (l *Lockout) GeneratePhase5a5bZkProof() (msg *models.Phase5A, err error) {
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

func (l *Lockout) ReceivePhase5A5BProof(msg *models.Phase5A, index int) (finish bool, err error) {
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
	e := CreateCommitmentWithUserDefinedRandomNess(inputhash.D, msg.Phase5ADecom1.BlindFactor)
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

//const
//decomVec:Di
//comVec: Ci
func phase5c(l *models.LocalSignature, decomVec []*models.Phase5ADecom1,

	vi *share.SPubKey,

) (*models.Phase5Com2, *models.Phase5DDecom2, error) {

	//从广播的commit(Ci,Di)得到vi,ai
	v := vi.Clone()
	for i := 0; i < len(decomVec); i++ {
		v.X, v.Y = share.PointAdd(v.X, v.Y, decomVec[i].Vi.X, decomVec[i].Vi.Y)
	}
	a := decomVec[0].Ai.Clone()
	for i := 1; i < len(decomVec); i++ {
		a.X, a.Y = share.PointAdd(a.X, a.Y, decomVec[i].Ai.X, decomVec[i].Ai.Y)
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
	com := CreateCommitmentWithUserDefinedRandomNess(inputhash.D, blindFactor)
	return &models.Phase5Com2{com},
		&models.Phase5DDecom2{
			Ui:          &share.SPubKey{uix, uiy}, //Ci
			Ti:          &share.SPubKey{tix, tiy}, //Di
			BlindFactor: blindFactor,
		},
		nil
}

func (l *Lockout) GeneratePhase5CProof() (msg *models.Phase5C, err error) {
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

func (l *Lockout) ReceivePhase5cProof(msg *models.Phase5C, index int) (finish bool, err error) {
	err = l.loadLockout()
	if err != nil {
		return
	}
	//验证5c的hash(ui,ti)=5c的Ci

	inputhash := proofs.CreateHashFromGE([]*share.SPubKey{msg.Phase5DDecom2.Ui, msg.Phase5DDecom2.Ti})
	inputhash.D = CreateCommitmentWithUserDefinedRandomNess(inputhash.D, msg.Phase5DDecom2.BlindFactor)
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

func (l *Lockout) Generate5dProof() (si share.SPrivKey, err error) {
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
	si, err = phase5d(l.L.LocalSignature, deCommitments2, commitments2, deCommitments1)
	if err != nil {
		return
	}
	l.L.Phase5D = make(map[int]share.SPrivKey)
	l.L.Phase5D[l.index] = si
	err = l.db.UpdateLockout(l.L)
	return
}

//const
//decom_vec2:	5c-Di
//com_vec2:		5c-Ci
//decom_vec1:	5a-Di
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

func (l *Lockout) RecevieSI(si share.SPrivKey, index int) (signature []byte, finish bool, err error) {
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

//const
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
	if share.BigInt2PrivateKey(gu1x).D.Cmp(r.D) == 0 {
		//return true
	}
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
