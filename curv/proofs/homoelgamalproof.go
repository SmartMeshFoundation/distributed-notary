package proofs

import (
	"math/big"

	"fmt"

	"encoding/hex"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/nkbai/goutils"
	"github.com/nkbai/log"
)

// HomoELGamalProof This is a proof of knowledge that a pair of group elements {D, E}
// form a valid homomorphic ElGamal encryption (”in the exponent”) using public key Y .
// (HEG is defined in B. Schoenmakers and P. Tuyls. Practical Two-Party Computation Based on the Conditional Gate)
// Specifically, the witness is ω = (x, r), the statement is δ = (G, H, Y, D, E).
// The relation R outputs 1 if D = xH+rY , E = rG (for the case of G=H this is ElGamal)
type HomoELGamalProof struct {
	T  *share.SPubKey
	A3 *share.SPubKey
	z1 share.SPrivKey
	z2 share.SPrivKey
}

// HomoElGamalWitness :
type HomoElGamalWitness struct {
	r share.SPrivKey
	x share.SPrivKey
}

// NewHomoElGamalWitness :
func NewHomoElGamalWitness(r, x share.SPrivKey) *HomoElGamalWitness {
	return &HomoElGamalWitness{r.Clone(), x.Clone()}
}

// HomoElGamalStatement :
type HomoElGamalStatement struct {
	G *share.SPubKey
	H *share.SPubKey
	Y *share.SPubKey
	D *share.SPubKey
	E *share.SPubKey
}

//CreateHomoELGamalProof const
func CreateHomoELGamalProof(w *HomoElGamalWitness, delta *HomoElGamalStatement) *HomoELGamalProof {
	s1 := share.RandomPrivateKey()
	s2 := share.RandomPrivateKey()
	A1x, A1y := S.ScalarMult(delta.H.X, delta.H.Y, s1.Bytes())
	A2x, A2y := S.ScalarMult(delta.Y.X, delta.Y.Y, s2.Bytes())
	A3x, A3y := S.ScalarMult(delta.G.X, delta.G.Y, s2.Bytes())
	tx, ty := share.PointAdd(A1x, A1y, A2x, A2y)
	e := CreateHashFromGE([]*share.SPubKey{{X: tx, Y: ty}, {X: A3x, Y: A3y}, delta.G, delta.H, delta.Y, delta.D, delta.E})
	z1 := s1.Clone()
	if w.x.D.Cmp(big.NewInt(0)) != 0 {
		t := e.Clone()
		t = share.ModMul(t, w.x)
		z1 = share.ModAdd(z1, t)
	}
	t := e.Clone()
	t = share.ModMul(t, w.r)
	z2 := s2.Clone()
	share.ModAdd(z2, t)
	return &HomoELGamalProof{
		T:  &share.SPubKey{X: tx, Y: ty},
		A3: &share.SPubKey{X: A3x, Y: A3y},
		z1: z1,
		z2: z2,
	}

}

//Verify : const 不会修改proof
func (proof *HomoELGamalProof) Verify(delta *HomoElGamalStatement) bool {
	e := CreateHashFromGE([]*share.SPubKey{proof.T, proof.A3, delta.G, delta.H, delta.Y, delta.D, delta.E})
	//z12=z1*H+z2*Y
	z12x, z12y := S.ScalarMult(delta.H.X, delta.H.Y, proof.z1.Bytes())
	x, y := S.ScalarMult(delta.Y.X, delta.Y.Y, proof.z2.Bytes())
	z12x, z12y = share.PointAdd(z12x, z12y, x, y)

	//T+e*D
	x, y = S.ScalarMult(delta.D.X, delta.D.Y, e.Bytes())
	tedx, tedy := share.PointAdd(x, y, proof.T.X, proof.T.Y)
	//z2g=G*z2
	z2gx, z2gy := S.ScalarMult(delta.G.X, delta.G.Y, proof.z2.Bytes())

	//A3+e*E
	x, y = S.ScalarMult(delta.E.X, delta.E.Y, e.Bytes())
	a3eex, a3eey := share.PointAdd(x, y, proof.A3.X, proof.A3.Y)

	if z12x.Cmp(tedx) == 0 && z12y.Cmp(tedy) == 0 &&
		z2gx.Cmp(a3eex) == 0 && z2gy.Cmp(a3eey) == 0 {
		return true
	}
	return false
}

// CreateHashFromGE :
func CreateHashFromGE(ge []*share.SPubKey) share.SPrivKey {
	var bs [][]byte
	for _, g := range ge {
		bs = append(bs, []byte{4})
		s := share.Xytostr(g.X, g.Y)
		log.Trace(fmt.Sprintf("s=%s", s))
		log.Trace(fmt.Sprintf("x=%s,y=%s", g.X.Text(16), g.Y.Text(16)))
		log.Trace(fmt.Sprintf("x=%s", hex.EncodeToString(g.X.Bytes())))
		log.Trace(fmt.Sprintf("write=%s", hex.EncodeToString(g.X.Bytes())))
		log.Trace(fmt.Sprintf("writey=%s", hex.EncodeToString(g.Y.Bytes())))
		bs = append(bs, g.X.Bytes())
		bs = append(bs, g.Y.Bytes())
		//bs = append(bs, b)
	}
	hash := utils.Sha256(bs...)
	result := new(big.Int).SetBytes(hash[:])
	return share.BigInt2PrivateKey(result)
}

/*func create_hash_from_ge(ge ...*ECPoint) *big.Int{
	var digest=sha256.New()
	for _,v:=range ge{

		tmp:=kgcenter.Get2Bytes(v.X,v.X)
		//digest=append(digest,tmp)
		digest.Write(tmp)
	}
	return new(big.Int).SetBytes(digest.Sum([]byte{}))
}*/

func pkToKeySlice() {

}
