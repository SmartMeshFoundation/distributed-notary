package messagetosign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

// BitcoinCancelPrepareLockoutTXDataName 用做消息传输时识别
const BitcoinCancelPrepareLockoutTXDataName = "BitcoinCancelPrepareLockoutTXData"

// BitcoinCancelPrepareLockoutTXData :
type BitcoinCancelPrepareLockoutTXData struct {
	tx              *wire.MsgTx
	lockScript      []byte
	notaryPublicKey *btcutil.AddressPubKey
	SCTokenAddress  common.Address `json:"sc_token_address"`
	SecretHash      common.Hash    `json:"secret_hash"`
	BytesToSign     []byte         `json:"bytes_to_sign"`
	Fee             int64          `json:"fee"`
}

// NewBitcoinCancelPrepareLockoutTXData :
func NewBitcoinCancelPrepareLockoutTXData(lockoutInfo *models.LockoutInfo, mcLockedPublicKey *btcutil.AddressPubKey, fee int64) *BitcoinCancelPrepareLockoutTXData {
	notaryAddress := mcLockedPublicKey.AddressPubKeyHash()
	// 4. 构造tx
	tx := wire.NewMsgTx(wire.TxVersion)
	tx.LockTime = uint32(lockoutInfo.MCExpiration)
	// txIn
	outpointTxHashInPrepareLockout, err := chainhash.NewHashFromStr(lockoutInfo.BTCPrepareLockoutTXHashHex)
	outpointInPrepareLockout := &wire.OutPoint{
		Hash:  *outpointTxHashInPrepareLockout,
		Index: lockoutInfo.BTCPrepareLockoutVout,
	}
	txIn := wire.NewTxIn(outpointInPrepareLockout, nil, nil)
	txIn.Sequence = 0
	tx.AddTxIn(txIn)
	// txOut
	pkScript, err := txscript.PayToAddrScript(notaryAddress)
	if err != nil {
		panic(err)
	}
	txOut := wire.NewTxOut(lockoutInfo.Amount.Int64()-fee, pkScript)
	tx.AddTxOut(txOut)
	// 5. 获取BytesToSign,
	lockScript := common.Hex2Bytes(lockoutInfo.BTCLockScriptHex)
	bytesToSign, err := txscript.CalcSignatureHash(lockScript, txscript.SigHashAll, tx, 0)
	if err != nil {
		panic(err)
	}
	data := &BitcoinCancelPrepareLockoutTXData{
		lockScript:      lockScript,
		tx:              tx,
		notaryPublicKey: mcLockedPublicKey,
		SCTokenAddress:  lockoutInfo.SCTokenAddress,
		SecretHash:      lockoutInfo.SecretHash,
		BytesToSign:     bytesToSign,
		Fee:             fee,
	}
	return data
}

// GetSignBytes : impl MessageToSign
func (d *BitcoinCancelPrepareLockoutTXData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *BitcoinCancelPrepareLockoutTXData) GetName() string {
	return BitcoinCancelPrepareLockoutTXDataName
}

// GetTransportBytes : impl MessageToSign
func (d *BitcoinCancelPrepareLockoutTXData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *BitcoinCancelPrepareLockoutTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to BitcoinCancelPrepareLockoutTXData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *BitcoinCancelPrepareLockoutTXData) VerifySignData(localLockoutInfo *models.LockoutInfo, mcLockedPublicKey *btcutil.AddressPubKey) (err error) {
	// 1. 因为取消交易有可能提前发,有可能超时后补偿发,所以不进行状态校验,能走到这里说明PrepareLockout交易已经发送过了,所以取消交易无脑发无所谓,没有任何风险
	// 2. 使用本地数据校验签名数据
	local := NewBitcoinCancelPrepareLockoutTXData(localLockoutInfo, mcLockedPublicKey, d.Fee)
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("BitcoinCancelPrepareLockoutTXData.VerifySignBytes() fail,maybe attack")
	}
	return
}

// BuildRawTransaction :
func (d *BitcoinCancelPrepareLockoutTXData) BuildRawTransaction(dsmSignature []byte) (rawTx *wire.MsgTx, err error) {
	if d.tx == nil {
		panic("never happen")
	}
	//sigScript, err := txscript.SignatureScript(tx, 0, lockScript, txscript.SigHashAll, notaryPrivateKey, true)
	//1. 获取sigScript
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(dsmSignature[:32])
	s.SetBytes(dsmSignature[32:64])
	signature := &btcec.Signature{R: r, S: s}
	sig := append(signature.Serialize(), byte(txscript.SigHashAll))
	sb := txscript.NewScriptBuilder()
	sb.AddData(sig)
	sb.AddData(d.notaryPublicKey.PubKey().SerializeCompressed())
	sigScript, err := sb.Script()
	if err != nil {
		log.Error(err.Error())
		return
	}
	// 4.拼装SignatureScript
	sb.Reset()
	sb.AddOps(sigScript)
	sb.AddOp(txscript.OP_FALSE)
	sb.AddData(d.lockScript)
	d.tx.TxIn[0].SignatureScript, err = sb.Script()
	if err != nil {
		log.Error(err.Error())
		return
	}
	rawTx = d.tx
	return
}
