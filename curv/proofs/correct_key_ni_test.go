package proofs

import (
	"crypto/rand"
	"testing"

	"github.com/nkbai/goutils"
)

func TestCreateNICorrectKeyProof(t *testing.T) {
	privKey, err := GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}
	t.Logf("pk=%s", utils.StringInterface(privKey, 5))
	//t.Log("Sigma length:",len(proofParams.Sigma))
	proofParams := CreateNICorrectKeyProof(privKey)
	for _, sigmaX := range proofParams.Sigma {
		t.Log(sigmaX)
	}
	if !proofParams.Verify(&privKey.PublicKey) {
		t.Error("not pass")
	}
}
