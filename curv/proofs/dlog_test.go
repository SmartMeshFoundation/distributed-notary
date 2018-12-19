package proofs

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/nkbai/goutils"
	"github.com/nkbai/log"
)

func TestProve(t *testing.T) {
	witness := big.NewInt(30)
	proof := Prove(share.BigInt2PrivateKey(witness))
	log.Trace(fmt.Sprintf("proof=%s", utils.StringInterface(proof, 7)))
	if !Verify(proof) {
		t.Error("should pass")
	}

}
