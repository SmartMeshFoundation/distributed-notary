package dnc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"

	"github.com/ethereum/go-ethereum/params"
)

func TestEther(t *testing.T) {
	fmt.Println(params.Ether)
	fmt.Println(uint64(params.Ether))
	fmt.Println(int64(params.Ether))
}

func TestVerify(t *testing.T) {

	s := `{"name":"User-SCPrepareLockin","request_id":"183e","sc_token_address":"0x88148d2f9e23769a143396d6124121a13d5a7c39","signer":"0x56771612bfae7fda173cb89579cc67876e34d6e7","secret_hash":"0x6de53701dcd7916f01f11687c3d0066ebe932ea59980e25f0191e03e693679a4","mc_user_address":"0x56771612bfae7fda173cb89579cc67876e34d6e7","sc_user_address":"0x56771612bfae7fda173cb89579cc67876e34d6e7","signature":"R7m0aUkoDmulgbDTQwzTrvllD4zxZ7N3vbiMuG1IPpR8PluLeHc3KjDD9cRDA+TIKWiBUndJjX7ua2RIW33vigA="}`

	var req userapi.SCPrepareLockinRequest
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		t.Error(err)
		return
	}
	b := req.VerifySign(&req)
	if !b {
		t.Error("should success")
	}
}
