package models

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/ethereum/go-ethereum/common"
)

// SignedKey :
type SignedKey struct {
	WI      share.SPrivKey
	Gwi     *share.SPubKey
	KI      share.SPrivKey
	GammaI  share.SPrivKey
	GGammaI *share.SPubKey
}

// SignMessgeModel :
type SignMessgeModel struct {
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

// Phase5A :
type Phase5A struct {
	Phase5Com1    *Phase5Com1
	Phase5ADecom1 *Phase5ADecom1
	Proof         *proofs.HomoELGamalProof
}

// Phase5C :
type Phase5C struct {
	*Phase5Com2
	*Phase5DDecom2
}

// Phase5Com1 :
type Phase5Com1 struct {
	Com *big.Int
}

// Phase5Com2 :
type Phase5Com2 struct {
	Com *big.Int
}

// Phase5ADecom1 :
type Phase5ADecom1 struct {
	Vi          *share.SPubKey
	Ai          *share.SPubKey
	Bi          *share.SPubKey
	BlindFactor *big.Int
}

// Phase5DDecom2 :
type Phase5DDecom2 struct {
	UI          *share.SPubKey
	Ti          *share.SPubKey
	BlindFactor *big.Int
}

// MessageBPhase2 :
type MessageBPhase2 struct {
	MessageBGamma *MessageB
	MessageBWi    *MessageB
}

// DeltaPhase3 :
type DeltaPhase3 struct {
	Delta share.SPrivKey
}

// SignBroadcastPhase1 广播给所有其他签名参与者
type SignBroadcastPhase1 struct {
	Com         *big.Int
	BlindFactor *big.Int
}

// MessageA 一对一定向传播给指定公证人,要一对一传递到所有此次签名参与者,不需要保存到数据库中,临时用即可
type MessageA struct {
	C []byte //paillier encion 文本
}

// MessageB 对于A的计算结果,计算完毕立即返回给指定公证人
type MessageB struct {
	C            []byte //pailler加密文本
	BProof       *proofs.DLogProof
	BetaTagProof *proofs.DLogProof
	Beta         share.SPrivKey //这个能否分开
}

// SignMessage :
type SignMessage struct {
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

// LocalSignature :
type LocalSignature struct {
	LI   share.SPrivKey
	RhoI share.SPrivKey //ρi
	R    *share.SPubKey
	SI   share.SPrivKey
	M    *big.Int
	Y    *share.SPubKey
}

func byte2SignBroadcastPhase1Map(data []byte) map[int]*SignBroadcastPhase1 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*SignBroadcastPhase1
	byte2Interface(data, &k)
	return k
}
func byte2MessageBPhase2(data []byte) map[int]*MessageBPhase2 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*MessageBPhase2
	byte2Interface(data, &k)
	return k
}
func byte2DeltaPhase3(data []byte) map[int]*DeltaPhase3 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*DeltaPhase3
	byte2Interface(data, &k)
	return k
}
func byte2Phase5A(data []byte) map[int]*Phase5A {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*Phase5A
	byte2Interface(data, &k)
	return k
}
func byte2Phase5C(data []byte) map[int]*Phase5C {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*Phase5C
	byte2Interface(data, &k)
	return k
}
func byte2SPrivKeyMap(data []byte) map[int]share.SPrivKey {
	if len(data) < 3 {
		return nil
	}
	var k map[int]share.SPrivKey
	byte2Interface(data, &k)
	return k
}
func byte2S(data []byte) []int {
	if len(data) < 3 {
		return nil
	}
	var k []int
	byte2Interface(data, &k)
	return k
}
func byte2SignedKey(data []byte) *SignedKey {
	if len(data) < 3 {
		return nil
	}
	var k SignedKey
	byte2Interface(data, &k)
	return &k
}
func byte2LocalSignature(data []byte) *LocalSignature {
	if len(data) < 3 {
		return nil
	}
	var k LocalSignature
	byte2Interface(data, &k)
	return &k
}
func fromSignMessageModel(p *SignMessgeModel) *SignMessage {
	p2 := &SignMessage{
		Phase1BroadCast: byte2SignBroadcastPhase1Map(p.Phase1BroadCast),
		Phase2MessageB:  byte2MessageBPhase2(p.Phase2MessageB),
		Phase3Delta:     byte2DeltaPhase3(p.Phase3Delta),
		Phase5A:         byte2Phase5A(p.Phase5A),
		Phase5C:         byte2Phase5C(p.Phase5C),
		Phase5D:         byte2SPrivKeyMap(p.Phase5D),
		AlphaGamma:      byte2SPrivKeyMap(p.AlphaGamma),
		AlphaWI:         byte2SPrivKeyMap(p.AlphaWI),
		Delta:           byte2SPrivKeyMap(p.Delta),
		Status:          p.Status,
		Sigma:           byte2SPrivKey(p.Sigma),
		R:               byte2SPubKey(p.R),
		SignedKey:       byte2SignedKey(p.SignedKey),
		S:               byte2S(p.S),
		LocalSignature:  byte2LocalSignature(p.LocalSignature),
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

func toSignMessageModle(p *SignMessage) *SignMessgeModel {
	p2 := &SignMessgeModel{
		Key:             p.Key[:],
		UsedPrivateKey:  p.UsedPrivateKey[:],
		Message:         p.Message,
		S:               interface2Byte(p.S, p.S == nil),
		SignedKey:       interface2Byte(p.SignedKey, p.SignedKey == nil),
		Phase1BroadCast: interface2Byte(p.Phase1BroadCast, p.Phase1BroadCast == nil),
		MessageA:        interface2Byte(p.MessageA, p.MessageA == nil),
		Phase2MessageB:  interface2Byte(p.Phase2MessageB, p.Phase2MessageB == nil),
		Phase3Delta:     interface2Byte(p.Phase3Delta, p.Phase3Delta == nil),
		Sigma:           interface2Byte(p.Sigma, false),
		R:               interface2Byte(p.R, p.R == nil),
		LocalSignature:  interface2Byte(p.LocalSignature, p.LocalSignature == nil),
		Phase5A:         interface2Byte(p.Phase5A, p.Phase5A == nil),
		Phase5C:         interface2Byte(p.Phase5C, p.Phase5C == nil),
		Phase5D:         interface2Byte(p.Phase5D, p.Phase5D == nil),
		AlphaGamma:      interface2Byte(p.AlphaGamma, p.AlphaGamma == nil),
		AlphaWI:         interface2Byte(p.AlphaWI, p.AlphaWI == nil),
		Delta:           interface2Byte(p.Delta, p.Delta == nil),
		Status:          p.Status,
	}
	return p2
}

// NewSignMessage :
func (db *DB) NewSignMessage(p *SignMessage) error {
	return db.Create(toSignMessageModle(p)).Error
}

// LoadSignMessage :
func (db *DB) LoadSignMessage(key common.Hash) (*SignMessage, error) {
	var pi SignMessgeModel
	err := db.Where(&SignMessgeModel{
		Key: key[:],
	}).First(&pi).Error
	if err != nil {
		return nil, err
	}
	return fromSignMessageModel(&pi), nil
}

//UpdateSignMessage 这个需要像lockin一样进行拆分,不要每次都存完整的
func (db *DB) UpdateSignMessage(l *SignMessage) error {
	return db.Save(toSignMessageModle(l)).Error
}
