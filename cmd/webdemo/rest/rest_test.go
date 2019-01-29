package rest

import (
	"encoding/json"
	"testing"

	"github.com/nkbai/goutils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

func TestBytes2Address(t *testing.T) {
	p := "042DFCD4CCB3ACFC01EB719DC02BF25C74FFE6349E9D9B59E6E20B10B5EEB658167A288A5CE72F80C74BD782DF047A66C91139BB359E52E3BA6B239397A5D7791C"
	addr := bytes2Address(common.Hex2Bytes(p))
	t.Logf("addr=%s", addr.String())
	privKey := common.Hex2Bytes("CD82212FDB2966FF8F39FAAD8EA6ED672B1CDAD967A3F1AB7BB1E59D6090BEB5")
	key, err := crypto.ToECDSA(privKey)
	if err != nil {
		t.Error(err)
		return
	}
	addr2 := crypto.PubkeyToAddress(key.PublicKey)
	assert.EqualValues(t, addr, addr2)
}

func TestSignature(t *testing.T) {
	s := `{"Chain":"side","Tx":{"nonce":"0x4","gasPrice":"0x430e23400","gas":"0xd331","to":"0x88148d2f9e23769a143396d6124121a13d5a7c39","value":"0x0","input":"0x7fd408d200000000000000000000000056771612bfae7fda173cb89579cc67876e34d6e7bb788c0ebf7269891687ad56693119ebb62abb85fddbb838b23e51218c061ef9","v":"0x0","r":"0x570F02C15E3ACCC3E3000E05654D4E28347ECCB8BBE77F16E18BD4CBCA879798","s":"0x8545AAC943ED0EE6F3AE3A0452CD7F68F7E691E309F626F012041932A82F8D71","hash":"0x5ea1fada758222613e1e6ac9b5b5aa3b46dc7de1a392d500f37d67157c72188f"},"TxHash":"0x23f709317cc5dbbfabaa11f68fc31ef355aa671cd51bd4e4d8c56051b430240e","Signer":"0x56771612bfae7fda173cb89579cc67876e34d6e7"}
`
	var req sendTxRequest
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		t.Error(err)
		return
	}
	trySignature(&req.Tx, req.TxHash, req.Signer)
}
func TestSignature2(t *testing.T) {
	s := `
{"Chain":"side","Tx":{"nonce":"0x1","gasPrice":"0x430e23400","gas":"0xc809","to":"0x88148d2f9e23769a143396d6124121a13d5a7c39","value":"0x0","input":"0x043d91800000000000000000000000003efd05b913a59c8795b546d4b499840d821921d0f226844fd61eb48809f5091f7cbe5a06fef68d5351e746ee4da0c526cdb7ba87","v":"0x0","r":"0xA975CF2B0391725FE8631C1C625E9CED37ACD6BD16DDF28F184FEF131C4C7E26","s":"0x32E3433E225C3655BD33FE61898BC9861CAC75210CDF1BE425617BFE4161229B","hash":"0xeff4c431a74077c33bf149d61efab1e2147fc1d3a44cda305702ffb7c117fca1"},"TxHash":"0xac52e5531afda841a6d30f3e5af79b183c8b98cd3219a526eb62325ae9edd54d","Signer":"0x3efd05b913a59c8795b546d4b499840d821921d0"}

`
	var req sendTxRequest
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		t.Error(err)
		return
	}
	MainChainEndpoint = "http://127.0.0.1:19888"
	SideChainEndpoint = "http://127.0.0.1:17888"
	tr, err := doSendTx(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tr=%s", utils.StringInterface(tr, 3))
}
func TestJson(t *testing.T) {
	type tstruct struct {
		Sig []byte
	}
	var rs = tstruct{}
	rs.Sig = common.Hex2Bytes("317D272DE8F3F504E9EDBCD905C207D404B699DE274B1D1D13D23496499FC6AFFAEB9388342146AC7BD285194B987332EF424B3C4FB7E23F74F37BA7B7FF542700")
	t.Logf("rs.sig=%s", utils.StringInterface(rs, 3))
	ss, err := json.Marshal(rs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ss=%s", string(ss))
}
