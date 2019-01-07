package service

import (
	"fmt"
	"testing"

	"encoding/json"

	"reflect"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

func TestSessionLogMsg(t *testing.T) {
	fmt.Println(SessionLogMsg(utils.NewRandomHash(), "123... %s %s", "aaa", "bbbb"))
}

func TestUnmarsha(t *testing.T) {
	type testStruct struct {
		A string
		B int
	}
	var r api.BaseResponse
	r.ErrorCode = api.ErrorCodeSuccess
	r.Data = &testStruct{
		A: "123",
		B: 5,
	}
	jsonStr := utils.ToJSONString(r)
	fmt.Println(jsonStr)
	var r2 api.BaseResponse
	json.Unmarshal([]byte(jsonStr), &r2)
	fmt.Println(utils.ToJSONString(r2))
	fmt.Println(r2.Data)
	fmt.Println(reflect.TypeOf(r2.Data))

	var tt testStruct
	buf, _ := json.Marshal(r2.Data)
	json.Unmarshal(buf, &tt)
	fmt.Println(utils.ToJSONString(tt))

	if _, ok := r2.Data.(*testStruct); ok {
		fmt.Println("======================", ok)
	}
}
