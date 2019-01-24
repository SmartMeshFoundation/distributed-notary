package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Spectrum/common"
	"github.com/SmartMeshFoundation/distributed-notary/testcode"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

type testR struct {
	BaseReq
	BaseReqWithResponse
	BaseReqWithSignature
	BaseReqWithSCToken
	BaseReqWithSessionID
	other string
}

func TestRequest(t *testing.T) {
	tr := testR{
		BaseReq:              NewBaseReq("APIName-testR"),
		BaseReqWithResponse:  NewBaseReqWithResponse(),
		BaseReqWithSignature: NewBaseReqWithSignature(utils.NewRandomAddress()),
		BaseReqWithSCToken:   NewBaseReqWithSCToken(utils.NewRandomAddress()),
		BaseReqWithSessionID: NewBaseReqWithSessionID(utils.NewRandomHash(), 1),
	}
	//fmt.Println(utils.ToJSONStringFormat(tr))

	c := make(chan Req, 1)
	c <- &tr
	d := <-c
	//fmt.Println(utils.ToJSONStringFormat(d))
	if _, ok := d.(*testR); ok {
		fmt.Println("testR")
	}

	if a, ok := d.(Req); ok {
		fmt.Println("Req", a.GetRequestName())
	}
	if a, ok := d.(*BaseReq); ok {
		fmt.Println("BaseReq ", a.GetRequestName())
	}

	if _, ok := d.(ReqWithResponse); ok {
		fmt.Println("ReqWithResponse ")
	}
	if _, ok := d.(ReqWithSCToken); ok {
		fmt.Println("ReqWithSCToken ")
	}
	if _, ok := d.(ReqWithSessionID); ok {
		fmt.Println("ReqWithSessionID ")
	}
	if _, ok := d.(ReqWithSignature); ok {
		fmt.Println("ReqWithSignature ")
	}
	fmt.Println("----------------------")
	d2 := d.(ReqWithResponse)
	go func() {
		for {
			resp := <-d2.GetResponseChan()
			fmt.Printf("receive response :\n%s\n", utils.ToJSONStringFormat(resp))
			if resp.ErrorCode == ErrorCodeSuccess {
				return
			}
		}
	}()
	d2.WriteErrorResponse(ErrorCodeException, "custom errorMsg")
	time.Sleep(time.Second)
	d2.WriteErrorResponse(ErrorCodePermissionDenied)
	time.Sleep(time.Second)
	d2.WriteSuccessResponse(struct {
		A string      `json:"a"`
		B interface{} `json:"b"`
	}{
		A: "aaaaa",
		B: 12567,
	})
	time.Sleep(time.Second)
}

func TestBaseReqWithSignature(t *testing.T) {
	type TestMsg struct {
		A common.Hash `json:"a"`
		B int         `json:"b"`
	}
	type TestRequest struct {
		BaseReqWithSignature
		Msg TestMsg
	}

	// 1. 构造request
	privateKey := testcode.GetTestPrivateKey1()
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)
	req := &TestRequest{
		BaseReqWithSignature: NewBaseReqWithSignature(sender),
	}
	fmt.Println("Before sign : \n", utils.ToJSONStringFormat(req))

	// 2. Sign
	req.Sign(privateKey)
	fmt.Println("After sign : \n", utils.ToJSONStringFormat(req))

	// 3. Verify Signature
	assert.EqualValues(t, true, req.VerifySign())

}
