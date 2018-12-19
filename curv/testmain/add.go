package main

import (
	"crypto/elliptic"
	"encoding/hex"
	"log"

	"math/big"

	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/nkbai/goutils"
)

/*
xxxxxkey1....!!! E5FBB16DA7903BB95C9F29AC34948156BDF2DAC6879586EEFAF88E6ED53632D6821A334B4F160F807177BD489939CF959D4D5F928B9ED07EC456E7ED27390627
,2018/12/12 16:35:29 key1 x=d63236d56e8ef8faee869587c6daf2bd56819434ac299f5cb93b90a76db1fbe5,y=27063927ede756c47ed09e8b925f4d9d95cf399948bd7771800f164f4b331a82
2018/12/12 16:35:29 keybin=04d63236d56e8ef8faee869587c6daf2bd56819434ac299f5cb93b90a76db1fbe527063927ede756c47ed09e8b925f4d9d95cf399948bd7771800f164f4b331a82
2018/12/12 16:35:29 add=04e07c019652e41c182603038c51ac09a4d512748e3f85a823a51f6b39cb2be6f3dc6ce03f6b3aa67cb6bf7e27d7becc22919d7850ec48c8464635f7d0e850a788
*/
func strtoxy(s string) (x, y *big.Int) {
	s1 := s[:64]
	s2 := s[64:]
	s1b, _ := hex.DecodeString(s1)
	s2b, _ := hex.DecodeString(s2)
	s1bc := make([]byte, len(s1b))
	s2bc := make([]byte, len(s1b))
	i := 0
	for j := len(s1b) - 1; j >= 0; j-- {
		s1bc[i] = s1b[j]
		s2bc[i] = s2b[j]
		i++
	}
	x = new(big.Int)
	x.SetBytes(s1bc)
	y = new(big.Int)
	y.SetBytes(s2bc)
	return
}
func xytostr(x, y *big.Int) string {
	x1 := utils.BigIntTo32Bytes(x)
	y1 := utils.BigIntTo32Bytes(y)
	x2 := make([]byte, len(x1))
	y2 := make([]byte, len(x1))
	i := 0
	for j := len(x1) - 1; j >= 0; j-- {
		x2[i] = x1[j]
		y2[i] = y1[j]
		i++
	}
	s := fmt.Sprintf("%s%s", hex.EncodeToString(x2), hex.EncodeToString(y2))
	return s
}
func main() {
	x, y := strtoxy("becc27fc0d115fc314e7574c979697e0bd9a559f8a17ad0953f6c7f0e284d4ac379c4fc62a26cc050f8e5f37a488d8ade9613b7671093864fdd9a7b0218933cc")
	s := xytostr(x, y)
	log.Printf("s=%s", s)
	return
	TestAddPoint()
}
func TestAddPoint() {
	x1, y1 := strtoxy("becc27fc0d115fc314e7574c979697e0bd9a559f8a17ad0953f6c7f0e284d4ac379c4fc62a26cc050f8e5f37a488d8ade9613b7671093864fdd9a7b0218933cc")
	x2, y2 := strtoxy("e59e705cb909acaba73cef8c4b8e775cd87cc0956e4045306d7ded41947f04c62ae5cf50a9316423e1d066326532f6f7eeea6c461984c5a339c33da6fe68e11a")
	x3, y3 := strtoxy("cb08a05d8917ecbb9178c1e50b984956ac5ac6706b24f45e1e41a958f8e74a771bc653c9c9741d30a8d6f9dfe2b12d3765b3b7d756dd4302195e6beb32a084d9")
	key1 := elliptic.Marshal(crypto.S256(), x1, y1)
	log.Printf("key1 x=%s,y=%s", x1.Text(16), y1.Text(16))
	log.Printf("keybin=%s", hex.EncodeToString(key1))
	//key1, _ := generateKeyPair()
	key2 := elliptic.Marshal(crypto.S256(), x2, y2)

	add, err := secp256k1.AddPoint(key1, key2)
	if err != nil {
		panic(err)
		return
	}
	if len(add) != len(key1) {
		panic("length error")
	}
	log.Printf("add=%s", hex.EncodeToString(add))
	x4, y4 := elliptic.Unmarshal(crypto.S256(), add)
	log.Printf("x4=%s,y4=%s", x4.Text(16), y4.Text(16))
	log.Printf("x3=%s,y3=%s", x3.Text(16), y3.Text(16))
}
