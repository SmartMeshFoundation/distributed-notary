package bitcoin

import (
	"testing"

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

func getTestData() (testSecret, testSecretHash, userPrivateKeyBytes, notaryPrivateKeyBytes []byte) {
	testSecret = common.HexToHash("0x630fbde9dd9e9a5ad33f01454e1c3a1a8821c78c9a886f61aa113cc5877b8166").Bytes()
	testSecretHash = utils.ShaSecret(testSecret[:]).Bytes()
	userPrivateKeyBytes, _ = hex.DecodeString("4d949ef677a600e449047eadb64b0686fcd24c4e820e3a3076f2cb5beb345c35")
	notaryPrivateKeyBytes, _ = hex.DecodeString("396b36331bfd0705f826a6df70f6dcb56dacab2e6e56c10a319c1b349d7bdb3e")
	return
}

func getTestRedeemTx(amount btcutil.Amount, pkScript []byte) *wire.MsgTx {
	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	txOut := wire.NewTxOut(int64(amount), pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)
	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)
	return redeemTx
}

func TestPrepareLockInScriptBuilder_GetSigScriptForNotary(t *testing.T) {
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

	// 验证脚本
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyCheckLockTimeVerify | txscript.ScriptVerifyCheckSequenceVerify
	vm, err := txscript.NewEngine(pkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println("NewEngine err : ", err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println("Execute err : ", err)
		return
	}
	fmt.Println("Transaction successfully signed")
}

func TestPrepareLockInScriptBuilder_GetSigScriptForUser(t *testing.T) {
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)

	// 获取测试数据
	_, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	userPrivateKey := PrivateKeyBytes2PrivateKey(userPrivateKeyBytes)
	amount := btcutil.Amount(1)
	builder := bs.GetPrepareLockInScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(570000))

	// 锁定脚本构造
	lockScript, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	redeemTx := getTestRedeemTx(builder.amount, pkScript)

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

	// 验证脚本
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyCheckLockTimeVerify
	vm, err := txscript.NewEngine(pkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println("NewEngine err : ", err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println("Execute err : ", err)
		return
	}
	fmt.Println("Transaction successfully signed")
	// fee
	fmt.Println("fee test")
	f, _ := bs.GetBtcRPCClient().EstimateFee(10)
	fmt.Println(btcutil.NewAmount(f * float64(redeemTx.SerializeSize()) / 1024))
}

func TestPrepareLockoutScriptBuilder_GetSigScriptForUser(t *testing.T) {

	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)

	// 获取测试数据
	secret, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	userPrivateKey := PrivateKeyBytes2PrivateKey(userPrivateKeyBytes)
	amount := btcutil.Amount(1)
	builder := bs.GetPrepareLockOutScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(570000))

	// 锁定脚本构造
	lockScript, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	redeemTx := getTestRedeemTx(builder.amount, pkScript)

	// 签名txout
	sigScript, err := txscript.SignatureScript(redeemTx, 0, lockScript, txscript.SigHashAll, userPrivateKey, true)

	// 构造SignatureScript
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddOps(builder.GetSigScriptForUser(secret))
	sb.AddData(lockScript)
	redeemTx.TxIn[0].SignatureScript, _ = sb.Script()
	fmt.Println(txscript.DisasmString(redeemTx.TxIn[0].SignatureScript))

	// 验证脚本
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyCheckLockTimeVerify
	vm, err := txscript.NewEngine(pkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println("NewEngine err : ", err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println("Execute err : ", err)
		return
	}
	fmt.Println("Transaction successfully signed")
	// fee
	fmt.Println("fee test")
	f, _ := bs.GetBtcRPCClient().EstimateFee(10)
	fmt.Println(btcutil.NewAmount(f * float64(redeemTx.SerializeSize()) / 1024))
}

func TestPrepareLockoutScriptBuilder_GetSigScriptForNotary(t *testing.T) {
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)

	// 获取测试数据
	_, secretHash, userPrivateKeyBytes, notaryPrivateKeyBytes := getTestData()
	userPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(userPrivateKeyBytes, &bs.net)
	notaryPublicKeyHash := PrivateKeyBytes2AddressPublicKeyHash(notaryPrivateKeyBytes, &bs.net)
	notaryPrivateKey := PrivateKeyBytes2PrivateKey(notaryPrivateKeyBytes)
	amount := btcutil.Amount(1)
	builder := bs.GetPrepareLockOutScriptBuilder(userPublicKeyHash, notaryPublicKeyHash, amount, secretHash, big.NewInt(100))

	// 锁定脚本构造
	lockScript, _, pkScript := builder.GetPKScript()

	// 模拟tx构造
	redeemTx := getTestRedeemTx(builder.amount, pkScript)

	// 签名txout
	redeemTx.TxIn[0].Sequence = 0
	//redeemTx.LockTime = uint32(builder.expiration.Int64())
	redeemTx.LockTime = 5
	sigScript, err := txscript.SignatureScript(redeemTx, 0, lockScript, txscript.SigHashAll, notaryPrivateKey, true)

	// 构造SignatureScript
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddOps(builder.GetSigScriptForNotary())
	sb.AddData(lockScript)
	redeemTx.TxIn[0].SignatureScript, _ = sb.Script()

	// 验证脚本
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyCheckLockTimeVerify | txscript.ScriptVerifyCheckSequenceVerify
	vm, err := txscript.NewEngine(pkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println("NewEngine err : ", err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println("Execute err : ", err)
		return
	}
	fmt.Println("Transaction successfully signed")
	// fee
	fmt.Println("fee test")
	f, _ := bs.GetBtcRPCClient().EstimateFee(10)
	fmt.Println(btcutil.NewAmount(f * float64(redeemTx.SerializeSize()) / 1024))
}

func TestParseScript(t *testing.T) {
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
