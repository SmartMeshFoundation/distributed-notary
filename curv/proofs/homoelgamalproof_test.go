package proofs

import (
	"testing"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
)

func TestCreateHomoELGamalProof(t *testing.T) {
	witness := &HomoElGamalWitness{
		x: share.RandomPrivateKey(),
		r: share.RandomPrivateKey(),
	}

	h := share.RandomPrivateKey()
	Hx, Hy := S.ScalarBaseMult(h.Bytes())
	y := share.RandomPrivateKey()
	Yx, Yy := S.ScalarBaseMult(y.Bytes())
	tx := new(big.Int).Set(Hx)
	ty := new(big.Int).Set(Hy)
	tx, ty = S.ScalarMult(tx, ty, witness.x.Bytes())
	tx2, ty2 := S.ScalarMult(Yx, Yy, witness.r.Bytes())

	Dx, Dy := share.PointAdd(tx, ty, tx2, ty2)

	Ex, Ey := S.ScalarBaseMult(witness.r.Bytes())
	delta := &HomoElGamalStatement{
		G: &share.SPubKey{X: S.Gx, Y: S.Gy},
		H: &share.SPubKey{X: Hx, Y: Hy},
		Y: &share.SPubKey{X: Yx, Y: Yy},
		D: &share.SPubKey{X: Dx, Y: Dy},
		E: &share.SPubKey{X: Ex, Y: Ey},
	}

	prove := CreateHomoELGamalProof(witness, delta)
	if !prove.Verify(delta) {
		t.Error("not pass")
	}
}

/*func TestCreateHomoELGamalProof(t *testing.T) {
	witness:=&HomoElGamalWitness{
		kgcenter.RandomFromZn(secp256k1.S256().N),
		kgcenter.RandomFromZn(secp256k1.S256().N),
	}
	G:=&ECPoint{secp256k1.S256().Gx,secp256k1.S256().Gy}
	h:=kgcenter.RandomFromZn(secp256k1.S256().N)
	Hx,Hy :=secp256k1.S256().ScalarMult(G.X,G.Y,h.Bytes())
	y:=kgcenter.RandomFromZn(secp256k1.S256().N)
	Yx,Yy:=secp256k1.S256().ScalarMult(G.X,G.Y,y.Bytes())

	D:=secp256k1.S256().Add(&kgcenter.Point{G.X,G.Y},
		kgcenter.PointMul(witness.r,Y))
	E:=kgcenter.PointMul(witness.r,&kgcenter.Point{G.X,G.Y})
	delta:=&HomoElGamalStatement{
		G,
		&ECPoint{H[0],H[1]},
		&ECPoint{Y[0],Y[1]},
		&ECPoint{D[0],D[1]},
		&ECPoint{E[0],E[1]},
	}
	prove:=CreateHomoELGamalProof(witness,delta)
	if !prove.Verify(delta){
		t.Error("not pass")
	}
}*/
/*
T=Secp256k1Point { purpose: "combine", ge: PublicKey(81baa5514553297c39fb59310ec0109dc421254e2b30a1fbecb5b40f7b5b4ff4aaba9fd8a469d3565b69112cc013e0132ffd5a749470579ac4c5973168b30088) }
A3=Secp256k1Point { purpose: "base_fe", ge: PublicKey(b595bf95b198816af105ef859399da5ed73eb206109b80f9b94ad163ae1d1ff344da876ce30322179e077868b6ccc025e7462289baa774f5edfa3bab7cb1b72d) }
e2=Secp256k1Scalar { purpose: "from_big_int", fe: SecretKey(96e9eba8e92484c4a7dc1c13887bed88375a5a2979eca144aec4b7a4ba17664b) }
test cryptographic_primitives::proofs::sigma_correct_homomrphic_elgamal_enc::tests::test_correct_general_homo_elgamal ... ok

*/
func TestCreateHashFromGE(t *testing.T) {
	x1, y1 := share.Strtoxy("c1bbd91326c1ce1881f8ca694d6c588e1fd7986683b973720085ac1a51610f88372e6568d68da50d015751f944c6ee8d7f61321559dee9eee51a7b08c852e2ff")
	x2, y2 := share.Strtoxy("f30bc95889e783507bc64fcb6e9a8f5e83dc3b0917bf38759bdc79331edd179a399c0dd1ccf2793d83652c31848613baa10fb23333ec59b4147a6f18cd9ae6eb")
	r := share.Str2bigint("eeaef14a20bfe6e2470b5ddfeef9ae3de244b47d2cd2cbb06785e861b035bca9")
	rs := CreateHashFromGE([]*share.SPubKey{
		{X: x1, Y: y1}, {X: x2, Y: y2},
	})
	if r.Cmp(rs.D) != 0 {
		t.Error("not equal")
	}
}
