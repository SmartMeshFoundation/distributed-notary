package mecdsa

import (
	"fmt"

	"crypto/rand"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/goutils"
)

type LockedIn struct {
	srv *NotaryService
	db  *models.DB
	Key common.Hash //此次lockedin唯一的key
}

func NewLockedIn(srv *NotaryService, db *models.DB, key common.Hash) *LockedIn {
	return &LockedIn{
		srv: srv,
		db:  db,
		Key: key,
	}
}

//生成自己的随机数,所有公证人的私钥片都会从这里面取走一部分
func createKeys() (share.SPrivKey, *proofs.PrivateKey) {
	ui := share.RandomPrivateKey()
	dk, err := proofs.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return ui, dk
}
func (l *LockedIn) GeneratePhase1PubKeyProof() (msg *models.KeyGenBroadcastMessage1, err error) {
	ui, dk := createKeys()
	p := &models.PrivateKeyInfo{
		Key:             l.Key,
		UI:              ui,
		PaillierPrivkey: dk,
		Status:          models.PrivateKeyNegotiateStatusInit,
	}
	err = l.db.NewPrivateKeyInfo(p)
	if err != nil {
		return
	}
	proof := proofs.Prove(ui)
	msg = &models.KeyGenBroadcastMessage1{proof}
	return
}

func (l *LockedIn) ReceivePhase1PubKeyProof(m *models.KeyGenBroadcastMessage1, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
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

//const
func CreateCommitmentWithUserDefinedRandomNess(message *big.Int, blindingFactor *big.Int) *big.Int {
	hash := utils.Sha256(message.Bytes(), blindingFactor.Bytes())
	b := new(big.Int)
	b.SetBytes(hash[:])
	return b
}

func (l *LockedIn) GeneratePhase2PaillierKeyProof() (msg *models.KeyGenBroadcastMessage2, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
	if err != nil {
		return
	}
	blindFactor := share.RandomBigInt()
	correctKeyProof := proofs.CreateNICorrectKeyProof(p.PaillierPrivkey)
	x, _ := share.S.ScalarBaseMult(p.UI.Bytes())
	com := CreateCommitmentWithUserDefinedRandomNess(x, blindFactor)
	msg = &models.KeyGenBroadcastMessage2{
		Com:             com,
		PaillierPubkey:  &p.PaillierPrivkey.PublicKey,
		CorrectKeyProof: correctKeyProof,
		BlindFactor:     blindFactor,
	}
	p.PaillierKeysProof2 = make(map[int]*models.KeyGenBroadcastMessage2)
	p.PaillierKeysProof2[l.srv.NotaryShareArg.Index] = msg
	err = l.db.KGUpdatePaillierKeysProof2(p)
	return
}

func (l *LockedIn) ReceivePhase2PaillierPubKeyProof(m *models.KeyGenBroadcastMessage2,
	index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
	if err != nil {
		return
	}

	_, ok := p.PaillierKeysProof2[index]
	if ok {
		err = fmt.Errorf("Paillier pubkey roof for %d already exist", index)
		return
	}
	if !m.CorrectKeyProof.Verify(m.PaillierPubkey) {
		err = fmt.Errorf("Paillier pubkey roof for %d verify not pass", index)
		return
	}
	if CreateCommitmentWithUserDefinedRandomNess(p.PubKeysProof1[index].Proof.PK.X, m.BlindFactor).Cmp(m.Com) != 0 {
		err = fmt.Errorf("blind factor error for %d", index)
		return
	}
	p.PaillierKeysProof2[index] = m
	err = l.db.KGUpdatePaillierKeysProof2(p)
	if err != nil {
		return
	}
	//vss,secretSahres:=feldman.Share(ThresholdCount,ShareCount,p.UI)
	//所有因素都凑齐了
	finish = len(p.PaillierKeysProof2) == params.ShareCount
	return
}

func (l *LockedIn) GeneratePhase3SecretShare() (msgs map[int]*models.KeyGenBroadcastMessage3, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
	if err != nil {
		return
	}
	vss, secretShares := feldman.Share(params.ThresholdCount, params.ShareCount, p.UI)
	msg := &models.KeyGenBroadcastMessage3{
		Vss:         vss,
		SecretShare: secretShares[l.srv.NotaryShareArg.Index],
		Index:       l.srv.NotaryShareArg.Index,
	}
	p.SecretShareMessage3 = make(map[int]*models.KeyGenBroadcastMessage3)
	p.SecretShareMessage3[l.srv.NotaryShareArg.Index] = msg
	err = l.db.KGUpdateSecretShareMessage3(p)
	if err != nil {
		return
	}
	msgs = make(map[int]*models.KeyGenBroadcastMessage3)
	for i := 0; i < params.ShareCount; i++ {
		if i == l.srv.NotaryShareArg.Index {
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

func (l *LockedIn) ReceivePhase3SecretShare(msg *models.KeyGenBroadcastMessage3, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
	if err != nil {
		return
	}
	_, ok := p.SecretShareMessage3[index]
	if ok {
		err = fmt.Errorf("secret shares for %d already received", index)
		return
	}
	if !EqualGE(p.PubKeysProof1[index].Proof.PK, msg.Vss.Commitments[0]) {
		err = fmt.Errorf("pubkey not match for %d", index)
		return
	}
	if !msg.Vss.ValidateShare(msg.SecretShare, l.srv.NotaryShareArg.Index+1) {
		err = fmt.Errorf("secret share error for %d", index)
		return
	}
	p.SecretShareMessage3[index] = msg
	err = l.db.KGUpdateSecretShareMessage3(p)
	if err != nil {
		return
	}
	finish = len(p.SecretShareMessage3) == params.ShareCount
	//if finish{
	//	var secretShares [] share.SPrivKey
	//	for i:=0;i<ShareCount;i++{
	//		secretShares=append(secretShares,p.SecretShareMessage3[i].SecretShare)
	//	}
	//	p.SecretShareMessage3[0].Vss.ValidateShare()
	//}
	//所有因素都凑齐了.可以计算我自己持有的私钥片
	if len(p.SecretShareMessage3) == params.ShareCount {
		p.XI = share.SPrivKey{new(big.Int)}
		for _, s := range p.SecretShareMessage3 {
			share.ModAdd(p.XI, s.SecretShare)
		}
	}
	err = l.db.KGUpdateXI(p)
	return
}
func (l *LockedIn) GeneratePhase4PubKeyProof() (msg *models.KeyGenBroadcastMessage4, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
	if err != nil {
		return
	}
	if p.XI.D == nil {
		panic("must wait for phase 3 finish")
	}
	proof := proofs.Prove(p.XI)
	msg = &models.KeyGenBroadcastMessage4{proof}
	p.LastPubkeyProof4 = make(map[int]*models.KeyGenBroadcastMessage4)
	p.LastPubkeyProof4[l.srv.NotaryShareArg.Index] = msg
	err = l.db.KGUpdateLastPubKeyProof4(p)
	return
}
func (l *LockedIn) ReceivePhase4VerifyTotalPubKey(msg *models.KeyGenBroadcastMessage4, index int) (finish bool, err error) {
	p, err := l.db.LoadPrivatedKeyInfo(l.Key)
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
func EqualGE(pubGB *share.SPubKey, mtaGB *share.SPubKey) bool {
	return pubGB.X.Cmp(mtaGB.X) == 0 && pubGB.Y.Cmp(mtaGB.Y) == 0
}
