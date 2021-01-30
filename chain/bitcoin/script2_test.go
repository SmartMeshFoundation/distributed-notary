package bitcoin

import (
	"bytes"
	"testing"

	"github.com/Toorop/go-bitcoind"

	"github.com/btcsuite/btcd/chaincfg"

	"math/big"

	"fmt"

	"encoding/hex"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	SERVER_HOST        = "127.0.0.1"
	SERVER_PORT        = 18443
	USER               = "test"
	PASSWD             = "test"
	USESSL             = false
	WALLET_PASSPHRASE  = "p1"
	WALLET_PASSPHRASE2 = "p2"
)

//user: moq7zgksGd1Qhoys2uF1rfvm3amfCvZmpg
//notary:mpjFg65NMHnH4aFQEKGToM1eu6icrbkQ7r
var userOriginTxHash = "d078db703e1848640ac9b605a77354a10b6eb8e315580dbc25b0222ffc5bbd7d"
var totalAmount btcutil.Amount = 10 * 1e8
var txFee btcutil.Amount = 1000
var txIndex uint32 = 0

func getTestData2() (testSecret, testSecretHash, userPrivateKeyBytes, notaryPrivateKeyBytes []byte) {
	testSecret = common.HexToHash("0x630fbde9dd9e9a5ad33f01454e1c3a1a8821c78c9a886f61aa113cc5877b8166").Bytes()
	testSecretHash = utils.ShaSecret(testSecret[:]).Bytes()
	userPrivateKeyBytes, _ = hex.DecodeString("4d949ef677a600e449047eadb64b0686fcd24c4e820e3a3076f2cb5beb345c35")
	notaryPrivateKeyBytes, _ = hex.DecodeString("396b36331bfd0705f826a6df70f6dcb56dacab2e6e56c10a319c1b349d7bdb3e")
	return
}
func getChainHash(s string) *chainhash.Hash {
	var h chainhash.Hash
	err := chainhash.Decode(&h, s)
	if err != nil {
		panic(err)
	}
	return &h
}
func getPrepareLockinTx(amount btcutil.Amount, pkScript []byte) *wire.MsgTx {
	prepareLockinTx := wire.NewMsgTx(wire.TxVersion)
	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut := wire.NewOutPoint(getChainHash(userOriginTxHash), txIndex)
	txIn := wire.NewTxIn(prevOut, nil, nil)
	prepareLockinTx.AddTxIn(txIn)
	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut := wire.NewTxOut(int64(amount), pkScript)
	prepareLockinTx.AddTxOut(txOut)
	return prepareLockinTx
}
func getPkScript(addr string) []byte {
	userAddress, err := btcutil.DecodeAddress(addr, &chaincfg.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	pkScript, err := txscript.PayToAddrScript(userAddress)
	if err != nil {
		panic(err)
	}
	return pkScript
}
func getTestRedeemTx2(amount btcutil.Amount, prepareLockinTxHash string, pkScript []byte) *wire.MsgTx {

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut := wire.NewOutPoint(getChainHash(prepareLockinTxHash), 0)
	txIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)
	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut := wire.NewTxOut(int64(amount), pkScript)
	redeemTx.AddTxOut(txOut)
	return redeemTx
}
func getHexTx(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		panic(err)
	}
	txHex := hex.EncodeToString(buf.Bytes())
	return txHex
}
func TestPrepareLockIn(t *testing.T) {
	testPrepareLockIn(t)
}
func testPrepareLockIn(t *testing.T) string {
	ast := assert.New(t)
	bc, err := bitcoind.New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	ast.Nil(err)
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	ast.Nil(err)

	// 获取测试数据
	_, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData2()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	t.Logf("userPublicKeyHash=%s", userPublicKeyHash)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	t.Logf("notaryPublicKeyHash=%s", notaryPublicKeyHash)
	userPrivateKey := PrivateKeyBytes2PrivateKey(userPrivateKeyBytes)
	amount := totalAmount - txFee
	builder := bs.GetPrepareLockInScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(400))

	// 锁定脚本构造
	_, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	prepareLockInTx := getPrepareLockinTx(builder.amount, pkScript)

	sigScript, err := txscript.SignatureScript(prepareLockInTx, 0, getPkScript(userPublicKeyHash.String()), txscript.SigHashAll, userPrivateKey, true)
	ast.Nil(err)
	prepareLockInTx.TxIn[0].SignatureScript = sigScript
	s1, _ := txscript.DisasmString(sigScript)
	s2, _ := txscript.DisasmString(getPkScript(userPublicKeyHash.String()))
	t.Logf("sigscript=%s", s1)
	t.Logf("pkscript=%s", s2)
	fmt.Println(txscript.DisasmString(prepareLockInTx.TxIn[0].SignatureScript))
	t.Logf("txHash=%s,\ntx=%s", prepareLockInTx.TxHash(), getHexTx(prepareLockInTx))
	txID, err := bc.SendRawTransaction(getHexTx(prepareLockInTx))
	ast.Nil(err)
	t.Logf("txID=%s", txID)
	err = bc.Generate(6)
	ast.Nil(err)
	return txID
}

//func TestPrepareLockInScriptBuilder_GetSigScriptForNotary2(t *testing.T) {
//	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
//	assert.Empty(t, err)
//
//	// 获取测试数据
//	secret, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData2()
//	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
//	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
//	notaryPrivateKey := PrivateKeyBytes2PrivateKey(notaryPrivateKeyBytes)
//	amount := totalAmount - txFee
//	builder := bs.GetPrepareLockInScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(100))
//
//	// 锁定脚本构造
//	lockScript, _, pkScript := builder.GetPKScript()
//
//	// 模拟tx构造
//	prepareLockinTx := getPrepareLockinTx(builder.amount, pkScript)
//
//	// 签名txin
//	sigScript, err := txscript.SignatureScript(prepareLockinTx, 0, lockScript, txscript.SigHashAll, notaryPrivateKey, true)
//
//	// 构造SignatureScript
//	sb := txscript.NewScriptBuilder()
//	sb.AddOps(sigScript)
//	sb.AddOps(builder.GetSigScriptForNotary(secret))
//	sb.AddData(lockScript)
//	redeemTx.TxIn[0].SignatureScript, _ = sb.Script()
//	fmt.Println(txscript.DisasmString(redeemTx.TxIn[0].SignatureScript))
//
//	// 验证脚本
//	flags := txscript.ScriptBip16 | txscript.ScriptVerifyCheckLockTimeVerify | txscript.ScriptVerifyCheckSequenceVerify
//	vm, err := txscript.NewEngine(pkScript, redeemTx, 0,
//		flags, nil, nil, -1)
//	if err != nil {
//		fmt.Println("NewEngine err : ", err)
//		return
//	}
//	if err := vm.Execute(); err != nil {
//		fmt.Println("Execute err : ", err)
//		return
//	}
//	fmt.Println("Transaction successfully signed")
//}

func TestPrepareLockInScriptBuilder_GetSigScriptForUser2(t *testing.T) {
	ast := assert.New(t)
	bc, err := bitcoind.New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	ast.Nil(err)
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	ast.Nil(err)
	prepareLockinTxHash := testPrepareLockIn(t)
	// 获取测试数据
	_, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData2()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	t.Logf("userPublicKeyHash=%s", userPublicKeyHash)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	t.Logf("notaryPublicKeyHash=%s", notaryPublicKeyHash)
	userPrivateKey := PrivateKeyBytes2PrivateKey(userPrivateKeyBytes)
	amount := totalAmount - txFee
	builder := bs.GetPrepareLockInScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(400))

	// 锁定脚本构造
	lockScript, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	redeemTx := getTestRedeemTx2(builder.amount, prepareLockinTxHash, pkScript)

	// 签名txout
	redeemTx.TxIn[0].Sequence = 0
	redeemTx.LockTime = uint32(builder.expiration.Int64())
	sigScript, err := txscript.SignatureScript(redeemTx, 0, lockScript, txscript.SigHashAll, userPrivateKey, true)

	// 构造SignatureScript
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddOps(builder.GetSigScriptForUser())
	sb.AddData(lockScript)
	redeemTx.TxIn[0].SignatureScript, _ = sb.Script()
	fmt.Println(txscript.DisasmString(redeemTx.TxIn[0].SignatureScript))

	txID, err := bc.SendRawTransaction(getHexTx(redeemTx))
	ast.Nil(err)
	t.Logf("txID=%s", txID)
	err = bc.Generate(6)
	ast.Nil(err)
}

func TestParseScript2(t *testing.T) {
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)

	// 获取测试数据
	secret, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	notaryPrivateKey := PrivateKeyBytes2PrivateKey(notaryPrivateKeyBytes)
	amount := btcutil.Amount(1)
	builder := bs.GetPrepareLockInScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(100))

	// 锁定脚本构造
	lockScript, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	redeemTx := getTestRedeemTx(builder.amount, pkScript)

	// 签名txin
	sigScript, err := txscript.SignatureScript(redeemTx, 0, lockScript, txscript.SigHashAll, notaryPrivateKey, true)

	// 构造SignatureScript
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddOps(builder.GetSigScriptForNotary(secret))
	sb.AddData(lockScript)
	redeemTx.TxIn[0].SignatureScript, _ = sb.Script()
	fmt.Println(txscript.DisasmString(redeemTx.TxIn[0].SignatureScript))
	// 解析
	info := &BTCOutpointRelevantInfo{
		SecretHash:    common.BytesToHash(secretHash),
		LockScriptHex: common.Bytes2Hex(lockScript),
	}
	fmt.Println("isCancelPrepareLockin : ", bs.isCancelPrepareLockin(redeemTx.TxIn[0], info))
	fmt.Println("isLockin : ", bs.isLockin(redeemTx.TxIn[0], info))
	// fee
	fmt.Println("fee test")
	f, _ := bs.GetBtcRPCClient().EstimateFee(10)
	fmt.Println(btcutil.NewAmount(f * float64(redeemTx.SerializeSize()) / 1024))

}
