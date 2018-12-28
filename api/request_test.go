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
	BaseRequest
	BaseNotaryRequest
	BaseCrossChainRequest
	other string
}

func TestRequest(t *testing.T) {
	tr := testR{
		BaseRequest: NewBaseRequest("APIName-testR"),
	}
	tr.SessionID = utils.NewRandomHash()
	fmt.Println(utils.ToJSONStringFormat(tr))

	c := make(chan Request, 1)
	c <- &tr
	d := <-c
	fmt.Println(utils.ToJSONStringFormat(d))
	if _, ok := d.(*testR); ok {
		fmt.Println("testR")
	}

	if a, ok := d.(Request); ok {
		fmt.Println("Request", a.GetRequestName())
	}
	if a, ok := d.(*BaseRequest); ok {
		fmt.Println("BaseRequest ", a.GetRequestName())
	}

	if a, ok := d.(NotaryRequest); ok {
		fmt.Println("NotaryRequest ", a.GetSessionID().String())
	}
	//if _, ok := d.(*BaseNotaryRequest); ok {
	//	fmt.Println("BaseNotaryRequest")
	//}
	if a, ok := d.(CrossChainRequest); ok {
		fmt.Println("CrossChainRequest ", a.GetSCTokenAddress().String())
	}
	//if _, ok := d.(*BaseCrossChainRequest); ok {
	//	fmt.Println("BaseCrossChainRequest")
	//}

	fmt.Println("----------------------")
	switch d.(type) {
	case CrossChainRequest:
		fmt.Println("deal CrossChainRequest")
	case NotaryRequest:
		fmt.Println("deal NotaryRequest")
	}
	fmt.Println("----------------------")
	go func() {
		for {
			resp := <-d.GetResponseChan()
			fmt.Printf("receive response :\n%s\n", utils.ToJSONStringFormat(resp))
			if resp.ErrorCode == ErrorCodeSuccess {
				return
			}
		}
	}()
	d.WriteErrorResponse(ErrorCodeException, "custom errorMsg")
	time.Sleep(time.Second)
	d.WriteErrorResponse(ErrorCodePermissionDenied)
	time.Sleep(time.Second)
	d.WriteSuccessResponse(struct {
		A string      `json:"a"`
		B interface{} `json:"b"`
	}{
		A: "aaaaa",
		B: 12567,
	})
	time.Sleep(time.Second)
	fmt.Println("----------------------")
	r1 := NewBaseRequest("r1")
	r2 := NewBaseRequest("r2")
	fmt.Println(utils.ToJSONStringFormat(r1))
	fmt.Println(utils.ToJSONStringFormat(r2))
}

func TestNotaryRequestSignature(t *testing.T) {
	type TestMsg struct {
		A common.Hash `json:"a"`
		B int         `json:"b"`
	}
	type TestRequest struct {
		BaseRequest
		BaseNotaryRequest
		Msg TestMsg
	}

	// 1. 构造request
	sessionID := utils.NewRandomHash()
	privateKey := testcode.GetTestPrivateKey1()
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)
	req := &TestRequest{
		BaseRequest:       NewBaseRequest("NotaryAPI-TestRequest"),
		BaseNotaryRequest: NewBaseNotaryRequest(sessionID, sender),
	}
	fmt.Println("Before sign : \n", utils.ToJSONStringFormat(req))

	// 2. Sign
	sig := req.Sign(privateKey)
	fmt.Println("After sign : \n", utils.ToJSONStringFormat(req))
	assert.EqualValues(t, sig, req.Signature)

	// 3. Verify Signature
	assert.EqualValues(t, true, req.VerifySignature())
}
