package messagetosign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

// BitcoinLockinTXDataName 用做消息传输时识别
const BitcoinLockinTXDataName = "BitcoinLockinTXData"

// BitcoinLockinTXData :
type BitcoinLockinTXData struct {
	tx              *wire.MsgTx
	lockScript      []byte
	secretScript    []byte
	notaryPublicKey *btcutil.AddressPubKey
	SCTokenAddress  common.Address `json:"sc_token_address"`
	SecretHash      common.Hash    `json:"secret_hash"`
	BytesToSign     []byte         `json:"bytes_to_sign"`
	Fee             int64          `json:"fee"`
}

// NewBitcoinLockinTXData :
func NewBitcoinLockinTXData(bs *bitcoin.BTCService, lockinInfo *models.LockinInfo, mcLockedPublicKey *btcutil.AddressPubKey, fee int64) (data *BitcoinLockinTXData, err error) {
	// 2. 获取双方地址
	userAddress, err := btcutil.DecodeAddress(lockinInfo.MCUserAddressHex, bs.GetNetParam())
	if err != nil {
		log.Error(err.Error())
		return
	}

	notaryAddress := mcLockedPublicKey.AddressPubKeyHash()
	// 3. 构造secretScript及lockScript
	builder := bs.GetPrepareLockInScriptBuilder(userAddress.(*btcutil.AddressPubKeyHash), notaryAddress, btcutil.Amount(lockinInfo.Amount.Int64()), lockinInfo.SecretHash[:], big.NewInt(int64(lockinInfo.MCExpiration)))
	lockScript, _, _ := builder.GetPKScript()
	// 4. 构造tx
	tx := wire.NewMsgTx(wire.TxVersion)
	// txIn
	prevOut := wire.NewOutPoint(lockinInfo.BTCPrepareLockinTXHash, lockinInfo.BTCPrepareLockinVout)
	txIn := wire.NewTxIn(prevOut, nil, nil)
	tx.AddTxIn(txIn)
	// txOut
	pkScript, err := txscript.PayToAddrScript(notaryAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	txOut := wire.NewTxOut(lockinInfo.Amount.Int64()-fee, pkScript)
	tx.AddTxOut(txOut)
	// 5. 获取BytesToSign,
	bytesToSign, err := txscript.CalcSignatureHash(lockScript, txscript.SigHashAll, tx, 0)
	if err != nil {
		log.Error(err.Error())
		return
	}
	data = &BitcoinLockinTXData{
		lockScript:      lockScript,
		secretScript:    builder.GetSigScriptForNotary(lockinInfo.Secret[:]),
		tx:              tx,
		notaryPublicKey: mcLockedPublicKey,
		SCTokenAddress:  lockinInfo.SCTokenAddress,
		SecretHash:      lockinInfo.SecretHash,
		BytesToSign:     bytesToSign,
		Fee:             fee,
	}
	return
}

// GetSignBytes : impl MessageToSign
func (d *BitcoinLockinTXData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *BitcoinLockinTXData) GetName() string {
	return BitcoinLockinTXDataName
}

// GetTransportBytes : impl MessageToSign
func (d *BitcoinLockinTXData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *BitcoinLockinTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to BitcoinLockinTXData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData 这里直接校验本地状态及SignBytes就好,因为SignBytes中已经包含了完整的tx信息
func (d *BitcoinLockinTXData) VerifySignData(bs *bitcoin.BTCService, localLockinInfo *models.LockinInfo, mcLockedPublicKey *btcutil.AddressPubKey) (err error) {
	// 1. 校验本地lockinInfo状态
	if localLockinInfo.SCLockStatus != models.LockStatusUnlock {
		err = fmt.Errorf("SCLockStatus wrong")
		return
	}
	if localLockinInfo.MCLockStatus != models.LockStatusLock {
		err = fmt.Errorf("MCLockStatus wrong")
		return
	}
	// 2. 使用本地数据获取MsgToSign
	local, err := NewBitcoinLockinTXData(bs, localLockinInfo, mcLockedPublicKey, d.Fee)
	if err != nil {
		return
	}
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("BitcoinLockinTXData.VerifySignBytes() fail,maybe attack")
	}
	return
}

// BuildRawTransaction 在这里组装可用的RawTransaction,需要对dms签名进行处理,在前后拼装btc要求的字节
func (d *BitcoinLockinTXData) BuildRawTransaction(dsmSignature []byte) (rawTx *wire.MsgTx, err error) {
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
	sb.AddOps(d.secretScript)
	sb.AddData(d.lockScript)
	d.tx.TxIn[0].SignatureScript, err = sb.Script()
	if err != nil {
		log.Error(err.Error())
		return
	}
	rawTx = d.tx
	return
}
