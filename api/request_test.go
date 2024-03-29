package api

import (
	"fmt"
	"testing"
	"time"

	"encoding/json"

	"github.com/SmartMeshFoundation/Spectrum/common"
	"github.com/SmartMeshFoundation/distributed-notary/testcode"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
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
		BaseReqWithSignature: NewBaseReqWithSignature(),
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
	req := &TestRequest{
		BaseReqWithSignature: NewBaseReqWithSignature(),
	}
	fmt.Println("Before sign : \n", utils.ToJSONStringFormat(req))

	// 2. Sign
	req.Sign(req, privateKey)
	fmt.Println("After sign : \n", utils.ToJSONStringFormat(req))

	// 3. Verify Signature
	assert.EqualValues(t, true, req.VerifySign(req))

}

func TestJsonRawMessage(t *testing.T) {

	type UpLoadSomething struct {
		Type   string
		Object interface{}
	}

	type File struct {
		FileName string
	}

	type Png struct {
		Wide  int
		Hight int
	}

	input := `
    {
        "type": "File",
        "object": {
            "filename": "for test"
        }
    }
    `
	var object json.RawMessage
	ss := UpLoadSomething{
		Object: &object,
	}
	if err := json.Unmarshal([]byte(input), &ss); err != nil {
		panic(err)
	}
	switch ss.Type {
	case "File":
		var f File
		if err := json.Unmarshal(object, &f); err != nil {
			panic(err)
		}
		println(f.FileName)
	case "Png":
		var p Png
		if err := json.Unmarshal(object, &p); err != nil {
			panic(err)
		}
		println(p.Wide)
	}
}
