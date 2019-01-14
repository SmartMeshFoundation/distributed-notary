package main

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestEther(t *testing.T) {
	fmt.Println(params.Ether)
	fmt.Println(uint64(params.Ether))
	fmt.Println(int64(params.Ether))
}
