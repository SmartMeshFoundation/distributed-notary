package models

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/ethereum/go-ethereum/common"
)

type SignedKey struct {
	WI      share.SPrivKey
	Gwi     *share.SPubKey
	KI      share.SPrivKey
	GammaI  share.SPrivKey
	GGammaI *share.SPubKey
}
type LockoutModel struct {
	Key             []byte `gorm:"primary_key"`
	UsedPrivateKey  []byte //使用哪个privatekey foreign key
	Message         []byte `gorm:"type:varchar(4096);"` //代签名消息
	S               []byte //此次签名指定的一组公证人,大于t即可
	SignedKey       []byte `gorm:"type:varchar(4096);"`
	Phase1BroadCast []byte `gorm:"type:varchar(4096);"`
	//phase2
	MessageA       []byte `gorm:"type:varchar(4096);"`
	Phase2MessageB []byte `gorm:"type:varchar(4096);"`

	//phase3
	Phase3Delta []byte `gorm:"type:varchar(4096);"`
	Sigma       []byte `gorm:"type:varchar(4096);"` //各自保管各自的sigma

	//phase4 r
	R              []byte `gorm:"type:varchar(4096);"`
	LocalSignature []byte `gorm:"type:varchar(4096);"`

	Phase5A []byte `gorm:"type:varchar(4096);"`
	Phase5C []byte `gorm:"type:varchar(4096);"`
	Phase5D []byte `gorm:"type:varchar(4096);"` //所有人的签名片集合

	AlphaGamma []byte `gorm:"type:varchar(4096);"` //中间结果,
	AlphaWI    []byte `gorm:"type:varchar(4096);"` //中间结果,不用广播

	Delta  []byte `gorm:"type:varchar(4096);"` //收集其他人发过来的delta 齐了以后求和,得到总delta
	Status int    //签名进行到哪个地步了? 失败应该重新协商,
}
type Phase5A struct {
	Phase5Com1    *Phase5Com1
	Phase5ADecom1 *Phase5ADecom1
	Proof         *proofs.HomoELGamalProof
}
type Phase5C struct {
	*Phase5Com2
	*Phase5DDecom2
}

type Phase5Com1 struct {
	Com *big.Int
}
type Phase5Com2 struct {
	Com *big.Int
}
type Phase5ADecom1 struct {
	Vi          *share.SPubKey
	Ai          *share.SPubKey
	Bi          *share.SPubKey
	BlindFactor *big.Int
}
type Phase5DDecom2 struct {
	Ui          *share.SPubKey
	Ti          *share.SPubKey
	BlindFactor *big.Int
}

type MessageBPhase2 struct {
	MessageBGamma *MessageB
	MessageBWi    *MessageB
}
type DeltaPhase3 struct {
	Delta share.SPrivKey
}

//广播给所有其他签名参与者
type SignBroadcastPhase1 struct {
	Com         *big.Int
	BlindFactor *big.Int
}

//一对一定向传播给指定公证人,要一对一传递到所有此次签名参与者,不需要保存到数据库中,临时用即可
type MessageA struct {
	C []byte //paillier encion 文本
}

//对于A的计算结果,计算完毕立即返回给指定公证人
type MessageB struct {
	C            []byte //pailler加密文本
	BProof       *proofs.DLogProof
	BetaTagProof *proofs.DLogProof
	Beta         share.SPrivKey //这个能否分开
}

type Lockout struct {
	Key             common.Hash
	UsedPrivateKey  common.Hash //使用哪个privatekey
	Message         []byte      //代签名消息
	S               []int       //此次签名指定的一组公证人,大于t即可
	SignedKey       *SignedKey
	Phase1BroadCast map[int]*SignBroadcastPhase1
	//phase2
	MessageA       *MessageA
	Phase2MessageB map[int]*MessageBPhase2

	//phase3
	Phase3Delta map[int]*DeltaPhase3
	Sigma       share.SPrivKey //各自保管各自的sigma

	//phase4 r
	R              *share.SPubKey
	LocalSignature *LocalSignature
	//phase5
	Phase5A map[int]*Phase5A
	Phase5C map[int]*Phase5C
	Phase5D map[int]share.SPrivKey //所有人的签名片集合

	AlphaGamma map[int]share.SPrivKey //中间结果,
	AlphaWI    map[int]share.SPrivKey //中间结果,不用广播

	Delta  map[int]share.SPrivKey //收集其他人发过来的delta 齐了以后求和,得到总delta
	Status int                    //签名进行到哪个地步了? 失败应该重新协商,

}
type LocalSignature struct {
	LI   share.SPrivKey
	RhoI share.SPrivKey //ρi
	R    *share.SPubKey
	SI   share.SPrivKey
	M    *big.Int
	Y    *share.SPubKey
}

func Byte2SignBroadcastPhase1Map(data []byte) map[int]*SignBroadcastPhase1 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*SignBroadcastPhase1
	Byte2Interface(data, &k)
	return k
}
func Byte2MessageBPhase2(data []byte) map[int]*MessageBPhase2 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*MessageBPhase2
	Byte2Interface(data, &k)
	return k
}
func Byte2DeltaPhase3(data []byte) map[int]*DeltaPhase3 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*DeltaPhase3
	Byte2Interface(data, &k)
	return k
}
func Byte2Phase5A(data []byte) map[int]*Phase5A {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*Phase5A
	Byte2Interface(data, &k)
	return k
}
func Byte2Phase5C(data []byte) map[int]*Phase5C {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*Phase5C
	Byte2Interface(data, &k)
	return k
}
func Byte2SPrivKeyMap(data []byte) map[int]share.SPrivKey {
	if len(data) < 3 {
		return nil
	}
	var k map[int]share.SPrivKey
	Byte2Interface(data, &k)
	return k
}
func Byte2S(data []byte) []int {
	if len(data) < 3 {
		return nil
	}
	var k []int
	Byte2Interface(data, &k)
	return k
}
func Byte2SignedKey(data []byte) *SignedKey {
	if len(data) < 3 {
		return nil
	}
	var k SignedKey
	Byte2Interface(data, &k)
	return &k
}
func Byte2LocalSignature(data []byte) *LocalSignature {
	if len(data) < 3 {
		return nil
	}
	var k LocalSignature
	Byte2Interface(data, &k)
	return &k
}
func fromLockoutModel(p *LockoutModel) *Lockout {
	p2 := &Lockout{
		Phase1BroadCast: Byte2SignBroadcastPhase1Map(p.Phase1BroadCast),
		Phase2MessageB:  Byte2MessageBPhase2(p.Phase2MessageB),
		Phase3Delta:     Byte2DeltaPhase3(p.Phase3Delta),
		Phase5A:         Byte2Phase5A(p.Phase5A),
		Phase5C:         Byte2Phase5C(p.Phase5C),
		Phase5D:         Byte2SPrivKeyMap(p.Phase5D),
		AlphaGamma:      Byte2SPrivKeyMap(p.AlphaGamma),
		AlphaWI:         Byte2SPrivKeyMap(p.AlphaWI),
		Delta:           Byte2SPrivKeyMap(p.Delta),
		Status:          p.Status,
		Sigma:           Byte2SPrivKey(p.Sigma),
		R:               Byte2SPubKey(p.R),
		SignedKey:       Byte2SignedKey(p.SignedKey),
		S:               Byte2S(p.S),
		LocalSignature:  Byte2LocalSignature(p.LocalSignature),
		//MessageA:        &MessageA{p.MessageA},
	}
	p2.Key.SetBytes(p.Key)
	p2.UsedPrivateKey.SetBytes(p.UsedPrivateKey)
	//空slice处理错误, 会是一个长度为的[]byte,内容是0
	if p.Message != nil && len(p.Message) > 1 {
		p2.Message = p.Message
	}
	if p.MessageA != nil && len(p.MessageA) > 1 {
		p2.MessageA = &MessageA{p.MessageA}
	}
	return p2
}

func toLockoutModle(p *Lockout) *LockoutModel {
	p2 := &LockoutModel{
		Key:             p.Key[:],
		UsedPrivateKey:  p.UsedPrivateKey[:],
		Message:         p.Message,
		S:               Interface2Byte(p.S, p.S == nil),
		SignedKey:       Interface2Byte(p.SignedKey, p.SignedKey == nil),
		Phase1BroadCast: Interface2Byte(p.Phase1BroadCast, p.Phase1BroadCast == nil),
		MessageA:        Interface2Byte(p.MessageA, p.MessageA == nil),
		Phase2MessageB:  Interface2Byte(p.Phase2MessageB, p.Phase2MessageB == nil),
		Phase3Delta:     Interface2Byte(p.Phase3Delta, p.Phase3Delta == nil),
		Sigma:           Interface2Byte(p.Sigma, false),
		R:               Interface2Byte(p.R, p.R == nil),
		LocalSignature:  Interface2Byte(p.LocalSignature, p.LocalSignature == nil),
		Phase5A:         Interface2Byte(p.Phase5A, p.Phase5A == nil),
		Phase5C:         Interface2Byte(p.Phase5C, p.Phase5C == nil),
		Phase5D:         Interface2Byte(p.Phase5D, p.Phase5D == nil),
		AlphaGamma:      Interface2Byte(p.AlphaGamma, p.AlphaGamma == nil),
		AlphaWI:         Interface2Byte(p.AlphaWI, p.AlphaWI == nil),
		Delta:           Interface2Byte(p.Delta, p.Delta == nil),
		Status:          p.Status,
	}
	return p2
}

func (db *DB) NewLockedout(p *Lockout) error {
	return db.Create(toLockoutModle(p)).Error
}

func (db *DB) LoadLockout(key common.Hash) (*Lockout, error) {
	var pi LockoutModel
	err := db.Where(&LockoutModel{
		Key: key[:],
	}).First(&pi).Error
	if err != nil {
		return nil, err
	}
	return fromLockoutModel(&pi), nil
}

func (db *DB) UpdateLockout(l *Lockout) error {
	return db.Save(toLockoutModle(l)).Error
}
