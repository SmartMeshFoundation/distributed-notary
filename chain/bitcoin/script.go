package bitcoin

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

// GetPrepareLockInScriptBuilder :
func (bs *BTCService) GetPrepareLockInScriptBuilder(userAddr, notaryAddr *btcutil.AddressPubKeyHash, amount *big.Int, lockSecretHashBytes []byte, expiration *big.Int) *PrepareLockInScriptBuilder {
	return &PrepareLockInScriptBuilder{
		userAddr:            userAddr,
		notaryAddr:          notaryAddr,
		amount:              amount,
		lockSecretHashBytes: lockSecretHashBytes,
		expiration:          expiration,
		net:                 bs.net,
	}
}

// PrepareLockInScriptBuilder :
type PrepareLockInScriptBuilder struct {
	userAddr            *btcutil.AddressPubKeyHash
	notaryAddr          *btcutil.AddressPubKeyHash
	amount              *big.Int
	lockSecretHashBytes []byte
	expiration          *big.Int
	net                 chaincfg.Params
}

/*
GetPKScript 生成锁定脚本及对应的PKScript
###
IF
	OP_HASH256 {{LockSecretHash}} OP_EQUALVERIFY OP_DUP OP_HASH160 {{分布式私钥对应的PublicKeyHash}} OP_EQUALVERIFY OP_CHECKSIG
ELSE
	{{过期的块号}} CHECKLOCKTIMEVERIFY DROP OP_DUP OP_HASH160 {{用户的PublicKeyHash}} OP_EQUALVERIFY OP_CHECKSIG
ENDIF
###

HASH160 {ADDR} EQUEAL
*/
func (b *PrepareLockInScriptBuilder) GetPKScript() (lockScript []byte, pkScript []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddOp(txscript.OP_IF)
	// 给公证人使用的锁定脚本
	sb.AddOp(txscript.OP_HASH256)
	sb.AddData(b.lockSecretHashBytes)
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_DUP)
	sb.AddOp(txscript.OP_HASH160)
	sb.AddData(b.notaryAddr.ScriptAddress())
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_CHECKSIG)

	sb.AddOp(txscript.OP_ELSE)
	// 给用户用于取消的脚本
	sb.AddInt64(b.expiration.Int64())
	sb.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
	sb.AddOp(txscript.OP_DROP)
	sb.AddOp(txscript.OP_DUP)
	sb.AddOp(txscript.OP_HASH160)
	sb.AddData(b.userAddr.ScriptAddress())
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_CHECKSIG)

	sb.AddOp(txscript.OP_ENDIF)
	lockScript, err := sb.Script()
	if err != nil {
		panic(err)
	}
	scriptHash, err := btcutil.NewAddressScriptHash(lockScript, &b.net)
	if err != nil {
		panic(err)
	}

	pkScript, err = txscript.PayToAddrScript(scriptHash)
	if err != nil {
		panic(err)
	}
	return
}

/*
GetSigScriptForNotary 获取公证人所使用的解锁脚本,不包含最前面的sig
###
sig {{分布式私钥对应的PublicKeyHash}} {{secretBytes}} OP_TRUE
###
*/
func (b *PrepareLockInScriptBuilder) GetSigScriptForNotary(secretBytes []byte) (script []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddData(secretBytes)
	sb.AddOp(txscript.OP_TRUE)
	script, err := sb.Script()
	if err != nil {
		panic(err)
	}
	return
}

/*
GetSigScriptForUser 获取用户过期取消的脚本
###
sig {{用户私钥对应的PublicKeyHash}} OP_FALSE
###
*/
func (b *PrepareLockInScriptBuilder) GetSigScriptForUser() (script []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddOp(txscript.OP_FALSE)
	script, err := sb.Script()
	if err != nil {
		panic(err)
	}
	return
}

// GetPrepareLockOutScriptBuilder :
func (bs *BTCService) GetPrepareLockOutScriptBuilder(userAddr, notaryAddr *btcutil.AddressPubKeyHash, amount *big.Int, lockSecretHashBytes []byte, expiration *big.Int) *PrepareLockOutScriptBuilder {
	return &PrepareLockOutScriptBuilder{
		userAddr:            userAddr,
		notaryAddr:          notaryAddr,
		amount:              amount,
		lockSecretHashBytes: lockSecretHashBytes,
		expiration:          expiration,
		net:                 bs.net,
	}
}

// PrepareLockOutScriptBuilder :
type PrepareLockOutScriptBuilder struct {
	userAddr            *btcutil.AddressPubKeyHash
	notaryAddr          *btcutil.AddressPubKeyHash
	amount              *big.Int
	lockSecretHashBytes []byte
	expiration          *big.Int
	net                 chaincfg.Params
}

/*
GetPKScript 生成锁定脚本及对应的PKScript
###
IF
	OP_HASH256 {{LockSecretHash}} OP_EQUALVERIFY OP_DUP OP_HASH160 {{用户的PublicKeyHash}} OP_EQUALVERIFY OP_CHECKSIG
ELSE
	{{过期的块号}} CHECKLOCKTIMEVERIFY DROP OP_DUP OP_HASH160 {{分布式私钥对应的PublicKeyHash}} OP_EQUALVERIFY OP_CHECKSIG
ENDIF
###

HASH160 {ADDR} EQUEAL
*/
func (b *PrepareLockOutScriptBuilder) GetPKScript() (lockScript []byte, pkScript []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddOp(txscript.OP_IF)
	// 给用户使用的锁定脚本
	sb.AddOp(txscript.OP_HASH256)
	sb.AddData(b.lockSecretHashBytes)
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_DUP)
	sb.AddOp(txscript.OP_HASH160)
	sb.AddData(b.userAddr.ScriptAddress())
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_CHECKSIG)

	sb.AddOp(txscript.OP_ELSE)
	// 给公证人用于取消的脚本
	sb.AddData(b.expiration.Bytes())
	sb.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
	sb.AddOp(txscript.OP_DROP)
	sb.AddOp(txscript.OP_DUP)
	sb.AddOp(txscript.OP_HASH160)
	sb.AddData(b.notaryAddr.ScriptAddress())
	sb.AddOp(txscript.OP_EQUALVERIFY)
	sb.AddOp(txscript.OP_CHECKSIG)

	sb.AddOp(txscript.OP_ENDIF)
	lockScript, err := sb.Script()
	if err != nil {
		panic(err)
	}
	scriptHash, err := btcutil.NewAddressScriptHash(lockScript, &b.net)
	if err != nil {
		panic(err)
	}

	pkScript, err = txscript.PayToAddrScript(scriptHash)
	if err != nil {
		panic(err)
	}
	return
}

/*
GetSigScriptForUser 获取用户过期取消的脚本
###
sig {{用户私钥对应的PublicKeyHash}} {{secret}} OP_TRUE
###
*/
func (b *PrepareLockOutScriptBuilder) GetSigScriptForUser(secret []byte) (script []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddData(secret)
	sb.AddOp(txscript.OP_TRUE)
	script, err := sb.Script()
	if err != nil {
		panic(err)
	}
	return
}

/*
GetSigScriptForNotary 获取公证人所使用的解锁脚本
###
sig {{分布式私钥对应的PublicKeyHash}} OP_FALSE
###
*/
func (b *PrepareLockOutScriptBuilder) GetSigScriptForNotary() (script []byte) {
	sb := txscript.NewScriptBuilder()
	sb.AddOp(txscript.OP_FALSE)
	script, err := sb.Script()
	if err != nil {
		panic(err)
	}
	return
}
