package feldman

import (
	"math/big"
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/nkbai/goutils"
	"github.com/nkbai/log"
	"github.com/stretchr/testify/assert"
)

func TestEvaluatePolynomial(t *testing.T) {
	cf := []share.SPrivKey{
		{D: share.Str2bigint("34930839620e77f1a7560698d20469b9e5f102f20980d204f73e6d37bb91f18c")},
		{D: share.Str2bigint("c5a8645d5d9c9f7362ccc677491a309c6e232573c7ed7a0bf1de65e25ac80772")},
		{D: share.Str2bigint("723981eb59fe890d67e536a2efc91f53d56a2bd37db7a0694bc46fbeb77b059c")},
	}
	//point := Str2bigint("0x0000000000000000000000000000000000000000000000000000000000000001")
	//res := Str2bigint("0x0bf422f5b4f7012edc2057ba2fca02eb89c6519a955ee611efc0afdf3825620a")
	secretShares := EvaluatePolynomial(cf, []int{1, 2, 3, 4, 5})
	t.Logf("secret_shares=%s", utils.StringInterface(secretShares, 5))
	//assert.EqualValues(t, res, secret_shares[0])
	/*
		15f8236f16cd73c8aff4bf2b6b071e208c7907767ae348d9827101f4303883a8
		b86a5219ac399bffb3b7cd83cb751a69e0e6795485dceda36ad96906e03505d9
		2bf2c59382898700252cad4aa2154f76e56ee7a66eb560bbf168878cbbe263bf
		70917ddc99bd34ca04535e7feee7bd450f700c3993fde29a95c31a9f63ad1fdc
		86467af4f1d4a55d512be123b1ec63d5a43b0a27466dd3039816c3b2075ef8ef
	*/
}
func TestScalarMult(t *testing.T) {
	a := share.Str2bigint("47626cae7657d2825645e60cf2d765f7470dfb55a6e2b65db1937fe6ad975d78")
	x, y := share.S.ScalarBaseMult(a.Bytes())
	s := share.Xytostr(x, y)
	t.Logf(s)
}
func TestShare(t *testing.T) {
	v, _ := Share(3, 5, share.BigInt2PrivateKey(big.NewInt(30)))
	t.Logf("v=%s", utils.StringInterface(v, 7))
}
func TestInvert(t *testing.T) {
	tests := []struct {
		a string
		r string
	}{
		{
			"0000000000000000000000000000000000000000000000000000000000000018",
			"f55555555555555555555555555555541d923e5d12a5998e97d44546f233fe89",
		},
		{
			"fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036413b",
			"2aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa74727a26728c1ab49ff8651778090ae0",
		},
		{
			"0000000000000000000000000000000000000000000000000000000000000004",
			"bfffffffffffffffffffffffffffffff0c0325ad0376782ccfddc6e99c28b0f1",
		},
	}
	for _, tt := range tests {
		a := share.Str2bigint(tt.a)
		ainv := share.Invert(a, share.S.N)
		if ainv.Cmp(share.Str2bigint(tt.r)) != 0 {
			t.Error("notequal")
		}
	}
}

func TestVerifiableSS_ValidateShare37(t *testing.T) {
	var seckey int64 = 99993993
	v, secretShares := Share(2, 7, share.BigInt2PrivateKey(big.NewInt(seckey)))
	log.Trace(fmt.Sprintf("v=%s", utils.StringInterface(v, 7)))
	s2 := secretShares[0:1]
	s2 = append(s2, secretShares[6])
	s2 = append(s2, secretShares[2])
	s2 = append(s2, secretShares[4])

	secretRescontructed := v.Reconstruct([]int{0, 6, 2, 4}, s2)
	if secretRescontructed.D.Cmp(big.NewInt(seckey)) != 0 {
		t.Error("reconstructed error")
		return
	}

	b := v.ValidateShare(secretShares[2], 3)
	assert.EqualValues(t, b, true)
	b = v.ValidateShare(secretShares[0], 1)
	assert.EqualValues(t, b, true)
	s := []int{0, 1, 3, 4, 6}
	l0 := v.MapShareToNewParams(0, s)
	l1 := v.MapShareToNewParams(1, s)
	l2 := v.MapShareToNewParams(3, s)
	l3 := v.MapShareToNewParams(4, s)
	l4 := v.MapShareToNewParams(6, s)
	log.Trace(fmt.Sprintf("l0=%s\n,l1=%s\n,l2=%s\n,l3=%s\n,l4=%s\n",
		l0, l1, l2,
		l3, l4,
	))
}

func TestVerifiableSS_ValidateShare(t *testing.T) {
	v, secretShares := Share(3, 5, share.BigInt2PrivateKey(big.NewInt(70)))
	log.Trace(fmt.Sprintf("v=%s", utils.StringInterface(v, 7)))
	s2 := secretShares[0:3]
	s2 = append(s2, secretShares[4])

	secretRescontructed := v.Reconstruct([]int{0, 1, 2, 4}, s2)
	if secretRescontructed.D.Cmp(big.NewInt(70)) != 0 {
		t.Error("reconstructed error")
		return
	}

	b := v.ValidateShare(secretShares[2], 3)
	assert.EqualValues(t, b, true)
	b = v.ValidateShare(secretShares[0], 1)
	assert.EqualValues(t, b, true)
	s := []int{0, 1, 2, 3, 4}
	l0 := v.MapShareToNewParams(0, s)
	l1 := v.MapShareToNewParams(1, s)
	l2 := v.MapShareToNewParams(2, s)
	l3 := v.MapShareToNewParams(3, s)
	l4 := v.MapShareToNewParams(4, s)
	log.Trace(fmt.Sprintf("l0=%s\n,l1=%s\n,l2=%s\n,l3=%s\n,l4=%s\n",
		l0, l1, l2, l3, l4,
	))
}

func TestVerifiableSS_ValidateShare2(t *testing.T) {
	v, secretShares := Share(2, 4, share.BigInt2PrivateKey(big.NewInt(70)))
	log.Trace(fmt.Sprintf("v=%s", utils.StringInterface(v, 7)))
	s2 := secretShares[0:3]
	//s2 = append(s2, secretShares[3])

	secretRescontructed := v.Reconstruct([]int{0, 1, 2}, s2)
	if secretRescontructed.D.Cmp(big.NewInt(70)) != 0 {
		t.Logf("secretRescontructed=%s", secretRescontructed)
		t.Error("reconstructed error")
		return
	}

	//secretRescontructed = v.Reconstruct([]int{0, 1, 2, 3, 4}, secretShares)
	//if secretRescontructed.Cmp(big.NewInt(70)) != 0 {
	//	t.Error("reconstructed error")
	//	return
	//}

	//b := v.ValidateShare(secretShares[2], 3)
	//assert.EqualValues(t, b, true)
	//b = v.ValidateShare(secretShares[0], 1)
	//assert.EqualValues(t, b, true)
	//s := []int{0, 1, 2, 3, 4}
	//l0 := v.MapShareToNewParams(0, s)
	//l1 := v.MapShareToNewParams(1, s)
	//l2 := v.MapShareToNewParams(2, s)
	//l3 := v.MapShareToNewParams(3, s)
	//l4 := v.MapShareToNewParams(4, s)
	//log.Trace(fmt.Sprintf("l0=%s\n,l1=%s\n,l2=%s\n,l3=%s\n,l4=%s\n",
	//	l0.Text(16), l1.Text(16), l2.Text(16),
	//	l3.Text(16), l4.Text(16),
	//))
}

func TestPointSub(t *testing.T) {
	type test struct {
		a string
		b string
		r string
	}
	tests := []test{
		{
			"dce405715d3ec4c9c70c58aae8eb0457e3ad67c8485e8cb9d22a44dad4f21b7939f296b6ba9d5735b72ad78789ab36b878c8c2c754b0c8cdeb899a0ad9720869",
			"20b8359b85ecdbd7d09bd894561999e84f439dea297d2dcce1e8e253f83b0fcdca9bd97cdb74e819939ff7ff3beb99b508cce98c0546e40d85c27bbde4678790",
			"88c3c8725ef2ce78ea2ef0f37c8276cb836df90f54299d6fad5361e8b0ee040e0499227d5d527f3a371e6658718247cafb8619637f6147b635b66d12c499f2e0",
		},
		{
			"721d36ade784af14bfa860a60d924424ea86c6d91183498102cfc2f2c3fdf571f13c35190e3fb75a9bf84ce369c931e8b8835c04ac73cbc236439775a3e17bb0",
			"487363c885290aaeaa20ed3c7ae9c6a6940ffcc4927c6b786448ae1e5eaa1f9afdb5fd37d2131e399d8fa7e66e61a6b382122d2b3efcab3dfb7f2131ffc263a7",
			"92eb16d71ef959e0a1545e4f3feee00bedd9e5a2d1b70488e06eb22c5ede2270199126c7aba395f7327bee0b9b9c4a0cab1886c7dd417119e9feb6fad54ebf29",
		},
		{
			"4a67285f62834bd70df60a7ddf65d8189a47115a7ecae3e9402ddc3568cd07495d475113bccd40ac5a97e4d4b30920c4a6be4cbd6bbb1c7378c7034512024554",
			"ee17a6e356f83102446d223799907097e36aa463106a72fe1dd5c1546437993bfdad201190ba07aacdff90fcffaf273952fdf1a068b14bef048270782cab150c",
			"931acea09f313af7f9ec12981c9483f4362db4aca23e4ce4c119dded29064d71ec5a42ad0852d156c4cccb81e467af50bfe2288d6d07fbfe05f77ef841d356a4",
		},
		{
			"8fd66d17e364ff01edfc4e5db6de1224e2ae3e15a2324a69c76d3afef655854c43f9d03dbb18c3de54c58f7a1995abcc3f23531bc61a5bc9abcc1018bc0b1ed2",
			"9eb54b723ed432e2f54cb0e25ee733fce2220af49b7feaf913fecef73b460a495d088786093fd81b114af142a9a56261813f8214b3404a64b5738548738fa70c",
			"d1ee2d44dc15e5dbc50d90ebb11a2f4c93c219eee668dbf58975b2ecb666090357e5c892ce26410d48faa6546491c9d66ff488da84bdf6c329ad529e0464807c",
		},
		{
			"9dc8964bae861270a31725a388f1addc34b777ab6d402691bfa37ae60f6ce7a018208506a7a0be092ebb9cf6c259b94660ad4666560adae5a90ebb2e3a8ac91c",
			"e5c3abafae67f8cafa05246df1ac9b0cefaa853a7c47ee87d9805abddf1b2dcee23ed4800ec4f4440fd0b752be97ee1a9f414d141fa57966efb4382b182a2469",
			"f88ebc09ec7d9fbadcad9d1a8b62fdbae26a7dc3d1ba6c077ac45c3b83843527cbe9a1352397ece20d90490d2a432b148fe99cf5fb0bf5a125354c21e69f965d",
		},
		{
			"4a74630709c49150d10f0e607670f2c73dc2865c8850c5d7514e386eed5cb299b34852ce4e6f4b192e8ffca93f502fe4e877d52065805e11c0899fefbf447a07",
			"b2cdbce035a29196c7dc1c79de6da01555645ada71174e7686111d00a55b8fa879ffb8a40c7bef590a944dd6a2d63d4472b30ac40b0f99d373350d0c8246d42d",
			"1c0946caaa5aa53c1b3dc50205e68e741aeb343e93fc0d8308dfe2e8c7e35e64c99c6d215710d8a2c6dcccd6bb8bc993b83ac243b44cb40d217d4854c4efae3a",
		},
	}
	for _, tt := range tests {
		ax, ay := share.Strtoxy(tt.a)
		bx, by := share.Strtoxy(tt.b)
		rx, ry := share.Strtoxy(tt.r)
		x, y := share.PointSub(ax, ay, bx, by)
		if x.Cmp(rx) != 0 || y.Cmp(ry) != 0 {
			t.Error("not equal")
		}
	}

}

func TestLsh(t *testing.T) {
	a := share.Str2bigint("28064132643846632695607237370921442439956604667885930229835793462081545041461")
	x := share.Str2bigint("72030963781072581219587713863143054310990281998109173424654097085162671284585")
	y := x.Lsh(x, 256)
	s := a.Add(a, y)
	t.Logf(s.Text(16))
	a = share.Str2bigint("64707283623258898139504461127142352865415264410956389544247533112393852326510")
	a.Lsh(a, 256)
	t.Logf("shl=%s", a.Text(10))

}

func TestPointAdd(t *testing.T) {
	x1, y1 := share.Strtoxy("4a74630709c49150d10f0e607670f2c73dc2865c8850c5d7514e386eed5cb299b34852ce4e6f4b192e8ffca93f502fe4e877d52065805e11c0899fefbf447a07")
	x2, y2 := share.Strtoxy("b2cdbce035a29196c7dc1c79de6da01555645ada71174e7686111d00a55b8fa879ffb8a40c7bef590a944dd6a2d63d4472b30ac40b0f99d373350d0c8246d42d")
	log.Trace(fmt.Sprintf("x1=%s,y1=%s", x1.Text(16), y1.Text(16)))
	log.Trace(fmt.Sprintf("x2=%s,y2=%s", x2.Text(16), y2.Text(16)))
	x, y := share.PointAdd(x1, y1, x2, y2)
	t.Logf("x=%s,y=%s", x.Text(16), y.Text(16))
	x, y = share.PointAdd(x2, y2, x1, y1)
	t.Logf("x=%s,y=%s", x.Text(16), y.Text(16))
}
