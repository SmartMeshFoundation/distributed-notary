package api

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

type testR struct {
	BaseRequest
	BaseNotaryRequest
	BaseCrossChainRequest
	other string
}

func TestRequest(t *testing.T) {
	tr := testR{}
	fmt.Println(utils.ToJsonStringFormat(tr), tr.GetRequestName())
	tr.Name = "name"
	tr.Key = utils.NewRandomHash()
	fmt.Println(utils.ToJsonStringFormat(tr), tr.GetRequestName())

	c := make(chan Request, 1)
	c <- &tr
	d := <-c
	fmt.Println(utils.ToJsonStringFormat(d))
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
		fmt.Println("NotaryRequest ", a.GetKey().String())
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
}
