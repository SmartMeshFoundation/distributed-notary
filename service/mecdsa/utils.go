package mecdsa

import (
	"crypto/rand"
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/proofs"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	utils2 "github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/goutils"
)

//phase1.1 生成自己的随机数,所有公证人的私钥片都会从这里面取走一部分
func createKeys() (share.SPrivKey, *proofs.PrivateKey) {
	ui := share.RandomPrivateKey()
	dk, err := proofs.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return ui, dk
}
func createCommitmentWithUserDefinedRandomNess(message *big.Int, blindingFactor *big.Int) *big.Int {
	hash := utils.Sha256(message.Bytes(), blindingFactor.Bytes())
	b := new(big.Int)
	b.SetBytes(hash[:])
	return b
}

// equalGE :
func equalGE(pubGB *share.SPubKey, mtaGB *share.SPubKey) bool {
	return pubGB.X.Cmp(mtaGB.X) == 0 && pubGB.Y.Cmp(mtaGB.Y) == 0
}

// sessionLogMsg :
func sessionLogMsg(sessionID common.Hash, formatter string, a ...interface{}) string {
	formatter = fmt.Sprintf("[SessionID=%s] %s", utils2.HPex(sessionID), formatter)
	if len(a) == 0 {
		return formatter
	}
	return fmt.Sprintf(formatter, a...)
}
