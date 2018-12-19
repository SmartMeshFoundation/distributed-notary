package models

import (
	"math/big"

	"errors"

	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/curv/feldman"
	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/ethereum/go-ethereum/common"
)

var errKeyLength = errors.New("key length error")

/*
每个公证人都有一个唯一的编号,0,1,2,3 ..., 不重复,
*/
/*
第一步,广播证明自己随机数对应公钥片
*/
type KeyGenBroadcastMessage1 struct {
	Proof *proofs.DLogProof
}

/*
第二步,广播证明自己此次私钥协商所用同态加密私钥
*/
type KeyGenBroadcastMessage2 struct {
	PaillierPubkey  *proofs.PublicKey         //paill ier 公钥
	Com             *big.Int                  //包含私钥片公钥信息的hash值
	CorrectKeyProof *proofs.NICorrectKeyProof //证明拥有一个paillier的私钥?
	BlindFactor     *big.Int                  //这有什么用呢?感觉不必要啊
}

/*
第三步,定向广播 Secret给指定的公证人
*/
type KeyGenBroadcastMessage3 struct {
	Vss         *feldman.VerifiableSS
	SecretShare share.SPrivKey
	Index       int
}
type KeyGenBroadcastMessage4 struct {
	Proof *proofs.DLogProof
}

const (
	PrivateKeyNegotiateStatusInit = iota
	PrivateKeyNegotiateStatusPubKey
	PrivateKeyNegotiateStatusPaillierPubKey
	PrivateKeyNegotiateStatusSecretShare
	PrivateKeyNegotiateStatusFinished
)

/*
lockedin 过程中互相之间协商的结果
*/
type PrivateKeyInfo struct {
	Key                 common.Hash
	PublicKeyY          *big.Int
	PublicKeyX          *big.Int                         // 此次协商生成的私钥对应的公钥 X,Y *big.Int
	UI                  share.SPrivKey                   //原始随机数,用于协商私钥片
	XI                  share.SPrivKey                   //自身私钥片
	PaillierPrivkey     *proofs.PrivateKey               //同态 私钥
	PubKeysProof1       map[int]*KeyGenBroadcastMessage1 //第一步广播的自身随机数对应公钥片证明信息
	PaillierKeysProof2  map[int]*KeyGenBroadcastMessage2 //第二步广播的同态加密公约证明信息
	SecretShareMessage3 map[int]*KeyGenBroadcastMessage3 //第三步,定向广播secretshare信息
	LastPubkeyProof4    map[int]*KeyGenBroadcastMessage4 //第四步,校验所有人收到的xi对应的pubkey,加总和一开始的总公钥是相同的.
	Status              int                              //Status
}
type PrivateKeyInfoModel struct {
	Key                 []byte `gorm:"primary_key"` //a random hash
	PublicKeyX          string // 此次协商生成的私钥对应的公钥 X,Y *big.Int
	PublicKeyY          string
	UI                  []byte `gorm:"type:varchar(128);"`  //原始随机数,用于协商私钥片
	XI                  []byte `gorm:"type:varchar(128);"`  //自身私钥片
	PaillierPrivkey     []byte `gorm:"type:varchar(1024);"` //同态 私钥
	PubKeysProof1       []byte `gorm:"type:varchar(4096);"` //第一步广播的自身随机数对应公钥片证明信息
	PaillierKeysProof2  []byte `gorm:"type:varchar(4096);"` //第二步广播的同态加密公约证明信息
	SecretShareMessage3 []byte `gorm:"type:varchar(4096);"` //第三步,定向广播secretshare信息
	LastPubkeyProof4    []byte `gorm:"type:varchar(4096);"` //第四步,校验所有人收到的xi对应的pubkey,加总和一开始的总公钥是相同的.
	Status              int
}

func Byte2Interface(data []byte, i interface{}) {
	buf := bytes.NewBuffer(data)
	d := gob.NewDecoder(buf)
	err := d.Decode(i)
	if err != nil {
		panic(fmt.Sprintf("decode err %s", err))
	}
	return
}
func Interface2Byte(i interface{}, isNil bool) []byte {
	if isNil {
		return nil
	}
	buf := new(bytes.Buffer)
	e := gob.NewEncoder(buf)
	err := e.Encode(i)
	if err != nil {
		panic(fmt.Sprintf("encode err %s ", err))
	}
	return buf.Bytes()
}
func Byte2KeyGenBroadcastMessage1Map(data []byte) map[int]*KeyGenBroadcastMessage1 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*KeyGenBroadcastMessage1
	Byte2Interface(data, &k)
	return k
}
func Byte2KeyGenBroadcastMessage2(data []byte) map[int]*KeyGenBroadcastMessage2 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*KeyGenBroadcastMessage2
	Byte2Interface(data, &k)
	return k
}
func Byte2KeyGenBroadcastMessage3(data []byte) map[int]*KeyGenBroadcastMessage3 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*KeyGenBroadcastMessage3
	Byte2Interface(data, &k)
	return k
}
func Byte2KeyGenBroadcastMessage4(data []byte) map[int]*KeyGenBroadcastMessage4 {
	if len(data) < 3 {
		return nil
	}
	var k map[int]*KeyGenBroadcastMessage4
	Byte2Interface(data, &k)
	return k
}

type tmpPrivateKey struct {
	P *big.Int
	Q *big.Int
}

func PaillierPrivateKey2Byte(sk *proofs.PrivateKey) []byte {
	p, q := sk.GetPQ()
	return Interface2Byte(&tmpPrivateKey{p, q}, false)

}
func Byte2PaillierPrivateKey(data []byte) *proofs.PrivateKey {
	if len(data) < 3 {
		return nil
	}
	k := &tmpPrivateKey{}
	Byte2Interface(data, k)
	return proofs.NewPrivateKey(k.P, k.Q)
}
func Byte2SPrivKey(data []byte) share.SPrivKey {
	if len(data) < 3 {
		return share.PrivKeyZero
	}
	k := &share.SPrivKey{}
	Byte2Interface(data, k)
	return *k
}
func Byte2SPubKey(data []byte) *share.SPubKey {
	if len(data) < 3 {
		return nil
	}
	k := &share.SPubKey{}
	Byte2Interface(data, k)
	return k
}
func strToBigInt(s string) *big.Int {
	if len(s) > 0 {
		i := new(big.Int)
		i.SetString(s, 0)
		return i
	}
	return nil
}
func bigIntToStr(i *big.Int) string {
	if i != nil {
		return i.String()
	}
	return ""
}
func fromPrivateKeyInfoModel(p *PrivateKeyInfoModel) *PrivateKeyInfo {
	p2 := &PrivateKeyInfo{
		PublicKeyX:          strToBigInt(p.PublicKeyX),
		PublicKeyY:          strToBigInt(p.PublicKeyY),
		UI:                  Byte2SPrivKey(p.UI),
		XI:                  Byte2SPrivKey(p.XI),
		PaillierPrivkey:     Byte2PaillierPrivateKey(p.PaillierPrivkey),
		PubKeysProof1:       Byte2KeyGenBroadcastMessage1Map(p.PubKeysProof1),
		PaillierKeysProof2:  Byte2KeyGenBroadcastMessage2(p.PaillierKeysProof2),
		SecretShareMessage3: Byte2KeyGenBroadcastMessage3(p.SecretShareMessage3),
		LastPubkeyProof4:    Byte2KeyGenBroadcastMessage4(p.LastPubkeyProof4),
		Status:              p.Status,
	}
	p2.Key.SetBytes(p.Key)
	return p2
}
func toPrivateKeyInfoModel(p *PrivateKeyInfo) *PrivateKeyInfoModel {
	p2 := &PrivateKeyInfoModel{
		Key:                 p.Key[:],
		PublicKeyX:          bigIntToStr(p.PublicKeyX),
		PublicKeyY:          bigIntToStr(p.PublicKeyY),
		UI:                  Interface2Byte(p.UI, false),
		XI:                  Interface2Byte(p.XI, false),
		PaillierPrivkey:     PaillierPrivateKey2Byte(p.PaillierPrivkey),
		PubKeysProof1:       Interface2Byte(p.PubKeysProof1, p.PubKeysProof1 == nil),
		PaillierKeysProof2:  Interface2Byte(p.PaillierKeysProof2, p.PaillierKeysProof2 == nil),
		SecretShareMessage3: Interface2Byte(p.SecretShareMessage3, p.SecretShareMessage3 == nil),
		Status:              p.Status,
	}
	return p2
}

func (db *DB) NewPrivateKeyInfo(p *PrivateKeyInfo) error {
	return db.Create(toPrivateKeyInfoModel(p)).Error
}

func (db *DB) LoadPrivatedKeyInfo(key common.Hash) (*PrivateKeyInfo, error) {
	var pi PrivateKeyInfoModel
	err := db.Where(&PrivateKeyInfoModel{
		Key: key[:],
	}).First(&pi).Error
	if err != nil {
		return nil, err
	}
	return fromPrivateKeyInfoModel(&pi), nil
}

//更新Paillier公钥信息
func (db *DB) KGUpdatePaillierKeysProof2(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		PaillierKeysProof2: Interface2Byte(p.PaillierKeysProof2, p.PaillierKeysProof2 == nil),
		Status:             PrivateKeyNegotiateStatusPaillierPubKey,
	}).Error
}

func (db *DB) KGUpdatePubKeysProof1(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		PubKeysProof1: Interface2Byte(p.PubKeysProof1, p.PubKeysProof1 == nil),
		Status:        PrivateKeyNegotiateStatusPubKey,
	}).Error
}
func (db *DB) KGUpdateTotalPubKey(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		PublicKeyX: bigIntToStr(p.PublicKeyX),
		PublicKeyY: bigIntToStr(p.PublicKeyY),
	}).Error
}
func (db *DB) KGUpdateSecretShareMessage3(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		SecretShareMessage3: Interface2Byte(p.SecretShareMessage3, p.SecretShareMessage3 == nil),
		Status:              PrivateKeyNegotiateStatusSecretShare,
	}).Error
}
func (db *DB) KGUpdateLastPubKeyProof4(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		LastPubkeyProof4: Interface2Byte(p.LastPubkeyProof4, p.LastPubkeyProof4 == nil),
	}).Error
}

func (db *DB) KGUpdateKeyGenStatus(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		Status: p.Status,
	}).Error
}

func (db *DB) KGUpdateXI(p *PrivateKeyInfo) error {
	return db.Model(&PrivateKeyInfoModel{
		Key: p.Key[:],
	}).Update(&PrivateKeyInfoModel{
		XI: Interface2Byte(p.XI, false),
	}).Error
}
