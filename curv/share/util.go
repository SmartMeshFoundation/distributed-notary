package share

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var S = secp256k1.S256()

var BigOne = big.NewInt(1)
var PrivKeyZero = BigInt2PrivateKey(big.NewInt(0))

//公钥
type SPubKey struct {
	X *big.Int
	Y *big.Int
}
type SPrivKey struct {
	D *big.Int
}

func (s SPrivKey) String() string {
	return s.D.Text(16)
}
func (s SPrivKey) Clone() SPrivKey {
	return SPrivKey{new(big.Int).Set(s.D)}
}
func (s SPrivKey) Bytes() []byte {
	return s.D.Bytes()
}
func NewGE(x, y *big.Int) *SPubKey {
	return &SPubKey{
		X: new(big.Int).Set(x),
		Y: new(big.Int).Set(y),
	}
}
func (g *SPubKey) Clone() *SPubKey {
	return &SPubKey{
		X: new(big.Int).Set(g.X),
		Y: new(big.Int).Set(g.Y),
	}
}
func (g *SPubKey) String() string {
	return Xytostr(g.X, g.Y)
}

// Invert calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.  fermatInverse
func Invert(k, N *big.Int) *big.Int {
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(k, nMinus2, N)
}
func InvertN(k SPrivKey) SPrivKey {
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(S.N, two)
	return SPrivKey{new(big.Int).Exp(k.D, nMinus2, S.N)}
}

func Str2bigint(s string) *big.Int {
	i := new(big.Int)
	i.SetString(s, 16)
	return i
}

func RandomPrivateKey() SPrivKey {
	key, _ := crypto.GenerateKey()
	return SPrivKey{key.D}
}
func RandomBigInt() *big.Int {
	r, err := rand.Int(rand.Reader, S.N)
	if err != nil {
		panic(fmt.Sprintf("err %s", err))
	}
	return r
}
func BigInt2PrivateKey(i *big.Int) SPrivKey {
	b := new(big.Int).Set(i)
	b.Mod(b, S.N)
	return SPrivKey{b}
}

func PointAdd(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	x,y=btcec.S256().Add(x1,y1,x2,y2)
	return
}
func PointSub(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	order := new(big.Int).Set(S.P)
	minusY := new(big.Int).Set(order)
	x = x2
	y = y2
	minusY = modSubInternal(minusY, y, order)
	return PointAdd(x1, y1, x2, minusY)

}
func Strtoxy(s string) (x, y *big.Int) {
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
func Xytostr(x, y *big.Int) string {
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

//s1=s1+s2 mod N
func ModAdd(s1, s2 SPrivKey) SPrivKey {
	s1.D.Mod(s1.D, S.N)
	t := new(big.Int).Mod(s2.D, S.N)
	s1.D.Add(s1.D, t)
	s1.D.Mod(s1.D, S.N)
	return s1
}

//s1=s1*s2 mod N
func ModMul(s1, s2 SPrivKey) SPrivKey {
	s1.D.Mod(s1.D, S.N)
	s2.D = new(big.Int).Mod(s2.D, S.N)
	s1.D.Mul(s1.D, s2.D)
	s1.D.Mod(s1.D, S.N)
	return s1
}

//s1=s1-s2 mod N
func ModSub(s1, s2 SPrivKey) SPrivKey {
	return SPrivKey{modSubInternal(s1.D, s2.D, S.N)}
}
func modSubInternal(s1, s2, modulus *big.Int) *big.Int {
	s1.Mod(s1, modulus)
	t := new(big.Int).Mod(s2, modulus)
	if s1.Cmp(t) >= 0 {
		s1.Sub(s1, t)
		return s1.Mod(s1, modulus)
	}
	big0 := big.NewInt(0)
	t = big0.Sub(big0, t)
	t.Add(t, modulus)
	s1.Add(s1, t)
	s1.Mod(s1, modulus)
	return s1
}

//s1=s1**s2 mod N
func ModPow(s1, s2 *big.Int) *big.Int {
	return s1.Exp(s1, s2, S.N)
}
