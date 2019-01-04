package userapi

import (
	"testing"
	"time"

	"fmt"

	"encoding/json"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

func TestUserAPI_Start(t *testing.T) {
	ua := NewUserAPI("127.0.0.1:8888")
	ua.SetTimeout(3 * time.Second)
	ua.Start(true)
}

func TestUserAPI_CreatePrivateKey(t *testing.T) {
	ua := NewUserAPI("127.0.0.1:8888")
	ua.SetTimeout(3 * time.Second)
	ua.createPrivateKey(nil, nil)
}

func TestRequest(t *testing.T) {
	r := &RegisterSCTokenRequest{
		BaseRequest: api.NewBaseRequest("Test"),
	}
	p := &struct {
		MainChainName string `json:"main_chain_name"`
		PrivateKeyID  string `json:"private_key_id"`
	}{
		MainChainName: "main-chain-name",
		PrivateKeyID:  utils.NewRandomHash().String(),
	}
	jsonStr := utils.ToJSONStringFormat(p)
	fmt.Println("before :", utils.ToJSONStringFormat(r))
	json.Unmarshal([]byte(jsonStr), &r)
	fmt.Println("after :", utils.ToJSONStringFormat(r))
}
