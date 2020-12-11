package messagetosign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"strings"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

// BitcoinPrepareLockoutTXDataName 用做消息传输时识别
const BitcoinPrepareLockoutTXDataName = "BitcoinPrepareLockoutTXData"

// BitcoinPrepareLockoutTXData :
type BitcoinPrepareLockoutTXData struct {
	originTx            *wire.MsgTx
	notaryAddressPubKey *btcutil.AddressPubKey
	lockScriptHex       string
	lockScriptHashHex   string
	//请求分发必须参数
	SCTokenAddress common.Address `json:"sc_token_address"`
	SecretHash     common.Hash    `json:"secret_hash"`
	//构造原始交易的必须数据
	UserRequest  *userapi.MCPrepareLockoutRequest `json:"user_request"`
	Fee          int64                            `json:"fee"`
	UTXOKeysStr  string                           `json:"utxo_keys_str"`
	TxInID       int                              `json:"tx_in_id"`      // 当前签名的txInID
	MCExpiration uint64                           `json:"mc_expiration"` // 由于公正之间的块号差距,这里以发起人的数据为准
	// 校验及签名数据
	OriginTXHash []byte `json:"origin_tx_hash"` // 原始交易的hash,校验数据用
	BytesToSign  []byte `json:"bytes_to_sign"`
}

// NewBitcoinPrepareLockoutTXData :
func NewBitcoinPrepareLockoutTXData(req *userapi.MCPrepareLockoutRequest, bs *bitcoin.BTCService, lockoutInfo *models.LockoutInfo, mcNotaryPublicKey *btcutil.AddressPubKey, db *models.DB, utxoKeysStr string, fee int64, indexToSign int) (data *BitcoinPrepareLockoutTXData, err error) {
	// 0. 获取本地utxos
	txHashs := strings.Split(utxoKeysStr, "-")
	var utxos []*models.BTCOutpoint
	for _, txHashStr := range txHashs {
		utxo, err2 := db.GetBTCOutpoint(txHashStr)
		if err2 != nil {
			err = err2
			log.Error(err.Error())
			return
		}
		if utxo.Status != models.BTCOutpointStatusUsable {
			err = fmt.Errorf("utxo %s can not use", txHashStr)
			return
		}
		utxos = append(utxos, utxo)
	}
	// 1. 获取双方地址
	userAddress := req.GetSignerBTCPublicKey(bs.GetNetParam()).AddressPubKeyHash()
	notaryAddress := mcNotaryPublicKey.AddressPubKeyHash()
	// 2. 构造原始交易
	tx := wire.NewMsgTx(wire.TxVersion)
	// txIn
	var totalAmount btcutil.Amount
	for _, utxo := range utxos {
		txIn := wire.NewTxIn(utxo.GetOutpoint(), nil, nil)
		tx.AddTxIn(txIn)
		totalAmount += utxo.Amount
	}
	// 找零txOut
	lockAmount := new(big.Int).Sub(lockoutInfo.Amount, lockoutInfo.CrossFee)
	backAmount := int64(totalAmount) - lockAmount.Int64() - fee
	if backAmount > 0 {
		pkScript, err2 := txscript.PayToAddrScript(notaryAddress)
		if err2 != nil {
			err = err2
			log.Error(err.Error())
			return
		}
		txOut4Notary := wire.NewTxOut(backAmount, pkScript)
		tx.AddTxOut(txOut4Notary)
	}
	// 锁定txOut
	builder := bs.GetPrepareLockOutScriptBuilder(userAddress, notaryAddress, btcutil.Amount(lockAmount.Int64()), lockoutInfo.SecretHash[:], big.NewInt(int64(lockoutInfo.MCExpiration)))
	lockScript, lockScriptAddr, _ := builder.GetPKScript()
	pkScript, err := txscript.PayToAddrScript(lockScriptAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	txOut4Lock := wire.NewTxOut(lockAmount.Int64(), pkScript)
	tx.AddTxOut(txOut4Lock)
	// 5. 生成BytesToSign,
	bytesToSign, err := txscript.CalcSignatureHash(utxos[indexToSign].GetPKScript(bs.GetNetParam()), txscript.SigHashAll, tx, indexToSign)
	if err != nil {
		log.Error(err.Error())
		return
	}
	originTXHash := tx.TxHash()
	data = &BitcoinPrepareLockoutTXData{
		originTx:            tx,
		notaryAddressPubKey: mcNotaryPublicKey,
		lockScriptHex:       common.Bytes2Hex(lockScript),
		lockScriptHashHex:   lockScriptAddr.String(),
		SCTokenAddress:      lockoutInfo.SCTokenAddress,
		SecretHash:          lockoutInfo.SecretHash,
		UserRequest:         req,
		Fee:                 fee,
		UTXOKeysStr:         utxoKeysStr,
		TxInID:              indexToSign,
		MCExpiration:        lockoutInfo.MCExpiration,
		OriginTXHash:        originTXHash.CloneBytes(),
		BytesToSign:         bytesToSign,
	}
	return
}

// GetSignBytes : impl MessageToSign
func (d *BitcoinPrepareLockoutTXData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *BitcoinPrepareLockoutTXData) GetName() string {
	return BitcoinPrepareLockoutTXDataName
}

// GetTransportBytes : impl MessageToSign
func (d *BitcoinPrepareLockoutTXData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *BitcoinPrepareLockoutTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to BitcoinLockinTXData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData 这里直接校验本地状态及SignBytes就好,因为SignBytes中已经包含了完整的tx信息
func (d *BitcoinPrepareLockoutTXData) VerifySignData(bs *bitcoin.BTCService, localLockoutInfo *models.LockoutInfo, mcNotaryPublicKey *btcutil.AddressPubKey, db *models.DB) (outpointToListen wire.OutPoint, err error) {
	// 1. 校验本地lockoutInfo状态
	if localLockoutInfo.SCLockStatus != models.LockStatusLock {
		err = fmt.Errorf("SCLockStatus wrong")
		return
	}
	if localLockoutInfo.MCLockStatus != models.LockStatusNone {
		err = fmt.Errorf("MCLockStatus wrong")
		return
	}
	// 2. 使用本地数据获取MsgToSign
	local, err := NewBitcoinPrepareLockoutTXData(d.UserRequest, bs, localLockoutInfo, mcNotaryPublicKey, db, d.UTXOKeysStr, d.Fee, d.TxInID)
	if err != nil {
		return
	}
	// 3. 校验用户请求
	//不校验了,因为jettrade这部分工作使用了不同的格式
	//if !d.UserRequest.VerifySign(d.UserRequest) {
	//	err = fmt.Errorf("signature in user request does't wrigt")
	//	return
	//}
	// 4. 校验原始交易
	if bytes.Compare(local.OriginTXHash, d.OriginTXHash) != 0 {
		err = fmt.Errorf("BitcoinPrepareLockoutTXData verify OriginTXHash fail,maybe attack")
		return
	}
	// 5. 校验SignBytes
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("BitcoinPrepareLockoutTXData verify SignBytes fail,maybe attack")
		return
	}
	// 返回outpoint及回填信息供注册使用
	outpointToListen = local.GetOriginTxCopy().TxIn[d.TxInID].PreviousOutPoint
	d.lockScriptHex = local.GetLockScriptHex()
	d.lockScriptHashHex = local.GetLockScriptHashHex()
	return
}

//GetOriginTxCopy :
func (d *BitcoinPrepareLockoutTXData) GetOriginTxCopy() *wire.MsgTx {
	return d.originTx.Copy()
}

// GetLockScriptHex :
func (d *BitcoinPrepareLockoutTXData) GetLockScriptHex() string {
	return d.lockScriptHex
}

// GetLockScriptHashHex :
func (d *BitcoinPrepareLockoutTXData) GetLockScriptHashHex() string {
	return d.lockScriptHashHex
}

// BuildBTCSignatureScript :
func (d *BitcoinPrepareLockoutTXData) BuildBTCSignatureScript(dsmSignature []byte) []byte {
	//1. 获取sigScript
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(dsmSignature[:32])
	s.SetBytes(dsmSignature[32:64])
	signature := &btcec.Signature{R: r, S: s}
	sig := append(signature.Serialize(), byte(txscript.SigHashAll))
	sb := txscript.NewScriptBuilder()
	sb.AddData(sig)
	sb.AddData(d.notaryAddressPubKey.PubKey().SerializeCompressed())
	signatureScript, err := sb.Script()
	if err != nil {
		panic(err)
	}
	return signatureScript
}
