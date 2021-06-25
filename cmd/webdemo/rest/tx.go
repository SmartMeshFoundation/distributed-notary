package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"reflect"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/ant0ine/go-json-rest/rest"

	putils "github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/utils"

	scontracts "github.com/SmartMeshFoundation/distributed-notary/chain/heco/contracts"
	mcontracts "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var hash2NetworkID = map[common.Hash]int{
	common.HexToHash("0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d"): 20180430, //spectrum主网
	common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"): 3,        //以太坊测试链
	common.HexToHash("0x38a88a9ddffe522df5c07585a7953f8c011c94327a494188bd0cc2410dc40a1a"): 8888,
	common.HexToHash("0xd011e2cc7f241996a074e2c48307df3971f5f1fe9e1f00cfa704791465d5efc3"): 3, //spectrum测试网
}

//网页端不好生成Tx
type TxRequest struct {
	From            string
	ContractAddress string                 //
	Method          string                 //prepareLockin,lockin,cancelLockin,prepareLockoutHTLC,lockout,cancleLockOut,queryLockin,queryLockout
	Arg             map[string]interface{} //根据各个函数调用不同,自己定义相应的参数
}

type TxResponse struct {
	TxHash common.Hash
	Tx     types.Transaction
	hasTx  bool
}

func generateTx(w rest.ResponseWriter, r *rest.Request) {
	req := &TxRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		errError(w, err)
		return
	}
	f := methodMap[req.Method]
	if f == nil {
		errString(w, "method not found")
	}
	from := common.HexToAddress(req.From)
	contract := common.HexToAddress(req.ContractAddress)
	tr, err := f(from, contract, req.Arg)
	if err != nil {
		errString(w, fmt.Sprintf("method %s err %s", req.Method, err))
		return
	}
	success(w, tr)
}

var methodMap = map[string]func(from, contract common.Address, m map[string]interface{}) (*TxResponse, error){
	"mprepareLockin":  mprepareLockin,
	"mcancelLockin":   mcancelLockin,
	"mlockout":        mlockout,
	"slockin":         slockin,
	"sprepareLockout": sprepareLockout,
	"scancelLockOut":  scancelLockOut,
}

func generateContractTx(endPoint string, from, contract common.Address, method string, args ...interface{}) (tr *TxResponse, err error) {
	conn, err := ethclient.Dial(endPoint)
	if err != nil {
		return
	}
	mc, err := mcontracts.NewLockedSpectrum(contract, conn)
	if err != nil {
		return
	}
	sc, err := scontracts.NewHecoToken(contract, conn)
	if err != nil {
		return
	}
	tr = &TxResponse{}
	transactor := &bind.TransactOpts{
		From: from,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != from {
				return nil, errors.New("not authorized to sign this account")
			}
			tr.TxHash.SetBytes(signer.Hash(tx).Bytes())
			tr.Tx = *tx
			tr.hasTx = true
			log.Trace(fmt.Sprintf("signer=%s", putils.StringInterface(signer, 5)))
			log.Trace(fmt.Sprintf("txhash=%s", tx.Hash().String()))
			return nil, errors.New("must fail")
		},
	}
	var err2 error
	switch method {
	case "mprepareLockin":
		if true {
			secretHash := args[0].(common.Hash)
			expiration := args[1].(*big.Int)
			value := args[2].(*big.Int)
			transactor.Value = value
			_, err2 = mc.PrepareLockin(transactor, secretHash, expiration)
		}
	case "mcancelLockin":
		if true {
			account := args[0].(common.Address)
			_, err2 = mc.CancelLockin(transactor, account)
		}
	case "mlockout":
		if true {
			account := args[0].(common.Address)
			secret := args[1].(common.Hash)
			_, err2 = mc.Lockout(transactor, account, secret)
		}
	case "slockin":
		if true {
			account := args[0].(common.Address)
			secret := args[1].(common.Hash)
			_, err2 = sc.Lockin(transactor, account, secret)
		}
	case "sprepareLockout":
		if true {
			secretHash := args[0].(common.Hash)
			expiration := args[1].(*big.Int)
			value := args[2].(*big.Int)
			_, err2 = sc.PrepareLockout(transactor, secretHash, expiration, value)
		}
	case "scancelLockOut":
		if true {
			account := args[0].(common.Address)
			_, err2 = sc.CancelLockOut(transactor, account)
		}
	}
	if !tr.hasTx {
		err = fmt.Errorf("try to contruct tx for %s ,err =%s", method, err2)
	}
	return
}
func mprepareLockin(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	secretHashStr, ok := m["SecretHash"]
	if !ok {
		err = fmt.Errorf("secret hash not found")
		return
	}
	SecretHash := common.HexToHash(secretHashStr.(string))
	if SecretHash == utils.EmptyHash {
		err = fmt.Errorf("Secret Hash error")
		return
	}
	ExpirationVal := m["Expiration"]
	log.Info(fmt.Sprintf("Expiration= %s", reflect.TypeOf(ExpirationVal)))
	Expiration := int64(ExpirationVal.(float64))
	if Expiration < int64(cfg.GetMinExpirationBlock4User(cfg.SMC.Name)) {
		err = fmt.Errorf("expiration error")
		return
	}
	value := big.NewInt(int64(m["Value"].(float64)))
	if value.Cmp(big.NewInt(0)) <= 0 {
		err = fmt.Errorf("value error")
		return
	}
	return generateContractTx(MainChainEndpoint, from, contract, "mprepareLockin", SecretHash, big.NewInt(Expiration), value)
}

func mcancelLockin(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	account := common.HexToAddress(m["Account"].(string))
	if account == utils.EmptyAddress {
		err = fmt.Errorf("Account err ")
		return
	}
	return generateContractTx(MainChainEndpoint, from, contract, "mcancelLockin", account)
}

func mlockout(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	account := common.HexToAddress(m["Account"].(string))
	if account == utils.EmptyAddress {
		err = fmt.Errorf("account error")
		return
	}
	secret := common.HexToHash(m["Secret"].(string))
	if secret == utils.EmptyHash {
		err = fmt.Errorf("secret must exist")
		return
	}
	return generateContractTx(MainChainEndpoint, from, contract, "mlockout", account, secret)
}

func slockin(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	account := common.HexToAddress(m["Account"].(string))
	secret := common.HexToHash(m["Secret"].(string))
	if account == utils.EmptyAddress {
		err = fmt.Errorf("account error")
		return
	}
	if secret == utils.EmptyHash {
		err = fmt.Errorf("secret error ")
		return
	}
	return generateContractTx(SideChainEndpoint, from, contract, "slockin", account, secret)
}
func sprepareLockout(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	SecretHash := common.HexToHash(m["SecretHash"].(string))
	expiration := int64(m["Expiration"].(float64))
	value := big.NewInt(int64(m["Value"].(float64)))
	if value.Cmp(big.NewInt(0)) <= 0 {
		err = fmt.Errorf("value err")
		return
	}
	if SecretHash == utils.EmptyHash {
		err = fmt.Errorf("empty secret hash")
		return
	}
	if expiration < int64(cfg.GetMinExpirationBlock4User(cfg.SMC.Name)) {
		err = fmt.Errorf("expiration error ")
		return
	}
	return generateContractTx(SideChainEndpoint, from, contract, "sprepareLockout", SecretHash, big.NewInt(expiration), value)
}

func scancelLockOut(from, contract common.Address, m map[string]interface{}) (tr *TxResponse, err error) {
	account := common.HexToAddress(m["Account"].(string))
	if account == utils.EmptyAddress {
		err = fmt.Errorf("account err")
		return
	}
	return generateContractTx(SideChainEndpoint, from, contract, "scancelLockout", account)
}

type sendTxRequest struct {
	Chain  string //main or side"
	Tx     txdata
	Signer common.Address
	TxHash common.Hash
}
type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

func (t txdata) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce hexutil.Uint64  `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     hexutil.Uint64  `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      hexutil.Bytes   `json:"input"    gencodec:"required"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.AccountNonce)
	enc.Price = (*hexutil.Big)(t.Price)
	enc.GasLimit = hexutil.Uint64(t.GasLimit)
	enc.Recipient = t.Recipient
	enc.Amount = (*hexutil.Big)(t.Amount)
	enc.Payload = t.Payload
	enc.V = (*hexutil.Big)(t.V)
	enc.R = (*hexutil.Big)(t.R)
	enc.S = (*hexutil.Big)(t.S)
	enc.Hash = t.Hash
	return json.Marshal(&enc)
}

func (t *txdata) UnmarshalJSON(input []byte) error {
	type txdata struct {
		AccountNonce *hexutil.Uint64 `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     *hexutil.Uint64 `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      *hexutil.Bytes  `json:"input"    gencodec:"required"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var dec txdata
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.AccountNonce == nil {
		return errors.New("missing required field 'nonce' for txdata")
	}
	t.AccountNonce = uint64(*dec.AccountNonce)
	if dec.Price == nil {
		return errors.New("missing required field 'gasPrice' for txdata")
	}
	t.Price = (*big.Int)(dec.Price)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gas' for txdata")
	}
	t.GasLimit = uint64(*dec.GasLimit)
	if dec.Recipient != nil {
		t.Recipient = dec.Recipient
	}
	if dec.Amount == nil {
		return errors.New("missing required field 'value' for txdata")
	}
	t.Amount = (*big.Int)(dec.Amount)
	if dec.Payload == nil {
		return errors.New("missing required field 'input' for txdata")
	}
	t.Payload = *dec.Payload
	if dec.V == nil {
		return errors.New("missing required field 'v' for txdata")
	}
	t.V = (*big.Int)(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' for txdata")
	}
	t.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 's' for txdata")
	}
	t.S = (*big.Int)(dec.S)
	if dec.Hash != nil {
		t.Hash = dec.Hash
	}
	return nil
}

func txData2Tx(t *txdata, txHash common.Hash, signer common.Address) *types.Transaction {
	trySignature(t, txHash, signer)
	buf, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	var tx types.Transaction
	err = json.Unmarshal(buf, &tx)
	if err != nil {
		panic(err)
	}
	return &tx

}
func getSig(t *txdata) []byte {
	r := putils.BigIntTo32Bytes(t.R)
	s := putils.BigIntTo32Bytes(t.S)
	var sig [65]byte
	copy(sig[:32], r)
	copy(sig[32:64], s)
	if t.V.Cmp(big.NewInt(0)) == 0 {
		sig[64] = 0
	} else {
		sig[64] = 1
	}
	return sig[:]
}

var halfN *big.Int

func init() {
	halfN, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	halfN.Div(halfN, big.NewInt(2))
}

/*
按照EIP155规范,必须halfN
*/
func halfS(t *txdata) {
	if t.S.Cmp(halfN) >= 0 {
		sig := make([]byte, 65)
		copy(sig[:32], putils.BigIntTo32Bytes(t.R))
		copy(sig[32:], putils.BigIntTo32Bytes(t.S))
		sig[64] = 0
		sig, err := secp256k1.SignatureNormalize(sig)
		if err != nil {
			log.Error(fmt.Sprintf("SignatureNormalize err %s\n,r=%s,s=%s", err, t.R, t.S))
			panic("normalize error")
		}
		t.S.SetBytes(sig[32:64])
	}
}
func trySignature(t *txdata, txHash common.Hash, signer common.Address) {
	halfS(t)
	r := putils.BigIntTo32Bytes(t.R)
	s := putils.BigIntTo32Bytes(t.S)
	var sig [65]byte
	copy(sig[:32], r)
	copy(sig[32:64], s)
	sig[64] = 0
	addr, err := utils.EcrecoverOnce(txHash, sig[:])
	if err == nil && addr == signer {
		t.V = big.NewInt(0)
		log.Info("0 is valid")
		return
	}
	sig[64] = 1
	addr, err = utils.EcrecoverOnce(txHash, sig[:])
	if err == nil && addr == signer {
		t.V = big.NewInt(1)
		log.Info("1 is valid")
		return
	}
	log.Error(fmt.Sprintf("already tried 0,1, but all wrong t=%s", putils.StringInterface(t, 3)))
}
func sendTx(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Error(fmt.Sprintf("sendTx err =%s", err))
		}
	}()
	var req sendTxRequest
	err = r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tr, err := doSendTx(req)
	if err != nil {
		errError(w, err)
	} else {
		success(w, tr)
	}
}

func doSendTx(req sendTxRequest) (tr *types.Receipt, err error) {
	log.Trace(fmt.Sprintf("req=%s", utils.ToJSONStringFormat(req)))
	endPoint := MainChainEndpoint
	if req.Chain == "side" {
		endPoint = SideChainEndpoint
	}
	conn, err := ethclient.Dial(endPoint)
	if err != nil {
		err = fmt.Errorf("eth dial %s", err)
		return
	}
	tx := txData2Tx(&req.Tx, req.TxHash, req.Signer)
	log.Trace(fmt.Sprintf("txData2Tx txhash=%s ", tx.Hash().String()))
	head, err := conn.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		err = fmt.Errorf("get header err %s", err)
		return
	}
	chainID := hash2NetworkID[head.Hash()]
	if chainID == 0 {
		if req.Chain == "side" {
			chainID = 7888
		} else {
			chainID = 9888
		}
	}
	log.Trace(fmt.Sprintf("chaind=%d,endpoint=%s,tx=%s", chainID, endPoint, utils.ToJSONStringFormat(&req.Tx)))
	signer := types.NewEIP155Signer(big.NewInt(int64(chainID)))
	tx, err = tx.WithSignature(signer, getSig(&req.Tx))
	if err != nil {
		err = fmt.Errorf("with signature err %s", err)
		return
	}
	log.Trace(fmt.Sprintf("after with signature,txhash=%s", tx.Hash().String()))
	sender, err := signer.Sender(tx)
	if err != nil {
		panic(err)
	}
	if sender != req.Signer {
		panic("not equal")
	}
	err = conn.SendTransaction(context.Background(), tx)
	if err != nil {
		err = fmt.Errorf("send transction err %s", err)
		return
	}
	tr, err = bind.WaitMined(context.Background(), conn, tx)
	if err != nil {
		err = fmt.Errorf("wait mined err %s", err)
		return
	}
	if tr.Status != 1 {
		err = fmt.Errorf("tx mined but failed")
		return
	}
	return
}
