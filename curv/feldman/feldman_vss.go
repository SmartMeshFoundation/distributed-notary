package feldman

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/ethereum/go-ethereum/crypto"
)

// ShamirSecretSharing :
type ShamirSecretSharing struct {
	Threshold  int
	ShareCount int
}

// VerifiableSS :
type VerifiableSS struct {
	Parameters  ShamirSecretSharing
	Commitments []*share.SPubKey //多项式的commit,第0个,也就是常数项,是secret对应的公钥
}

// Share :
func Share(t, n int, secret share.SPrivKey) (*VerifiableSS, []share.SPrivKey) {
	poly := SamplePolynomial(t, secret)
	var index []int
	for i := 1; i <= n; i++ {
		index = append(index, i)
	}
	secretShares := EvaluatePolynomial(poly, index)
	//log.Trace(fmt.Sprintf("secretShares=%s", secretShares))
	var commitments []*share.SPubKey
	for _, p := range poly {
		x, y := share.S.ScalarBaseMult(p.D.Bytes())
		commitments = append(commitments, &share.SPubKey{X: x, Y: y})
	}

	return &VerifiableSS{
		Parameters: ShamirSecretSharing{
			Threshold:  t,
			ShareCount: n,
		},
		Commitments: commitments,
	}, secretShares
}

// SamplePolynomial :
func SamplePolynomial(t int, coef0 share.SPrivKey) []share.SPrivKey {
	var bs []share.SPrivKey
	bs = append(bs, coef0)
	for i := 0; i < t; i++ {
		k, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		bs = append(bs, share.SPrivKey{D: k.D})
	}
	//bs = []*big.Int{
	//	coef0,
	//	big.NewInt(1),
	//	big.NewInt(2),
	//	big.NewInt(3),
	//}
	return bs
}

// EvaluatePolynomial :
func EvaluatePolynomial(coefficients []share.SPrivKey, index []int) []share.SPrivKey {
	var bs []share.SPrivKey
	for i := 0; i < len(index); i++ {
		point := share.BigInt2PrivateKey(big.NewInt(int64(index[i])))
		//log.Trace(fmt.Sprintf("point=%s", point))
		//log.Trace(fmt.Sprintf("coefficients=%s", utils.StringInterface(coefficients, 3)))
		sum := coefficients[len(coefficients)-1].Clone()
		for j := len(coefficients) - 2; j >= 0; j-- {
			//log.Trace(fmt.Sprintf("sum=%s,coef=%s", sum, coefficients[j]))
			share.ModMul(sum, point)
			share.ModAdd(sum, coefficients[j])
		}
		bs = append(bs, sum)
	}
	return bs
}

//ValidateShare : 验证我share出去的secret share,和index一一对应关系,index是下标加1 const
func (v *VerifiableSS) ValidateShare(secretShare share.SPrivKey, index int) bool {
	x, y := share.S.ScalarBaseMult(secretShare.D.Bytes())
	ssPoint := &share.SPubKey{X: x, Y: y}
	indexFe := big.NewInt(int64(index))
	indexFe = indexFe.Mod(indexFe, share.S.N)
	l := len(v.Commitments)
	//log.Trace(fmt.Sprintf("indexfe=%s", indexFe))
	head := v.Commitments[l-1].Clone()
	for j := l - 2; j >= 0; j-- {
		c := v.Commitments[j]
		//log.Trace(fmt.Sprintf("acc=%s,x=%s", Xytostr(head.X, head.Y), Xytostr(c.X, c.Y)))
		x, y = share.S.ScalarMult(head.X, head.Y, indexFe.Bytes())
		//log.Trace(fmt.Sprintf("t=%s", share.Xytostr(x, y)))
		x, y = share.PointAdd(x, y, c.X, c.Y)
		//log.Trace(fmt.Sprintf("x1=%s,y1=%s", x.Text(16), y.Text(16)))
		//x, y = S.Add(head.X, head.Y, x, y)
		//log.Trace(fmt.Sprintf("after add %s", Xytostr(x, y)))
		head = &share.SPubKey{X: x, Y: y}
	}
	//log.Trace(fmt.Sprintf("sspoint=%s,commit_to_point=%s", utils.StringInterface(ssPoint, 3), utils.StringInterface(head, 3)))
	return ssPoint.X.Cmp(head.X) == 0 && ssPoint.Y.Cmp(head.Y) == 0
}

//MapShareToNewParams 给定的参数,不会发生变化,这是系数? 给相同的index和s,结果返回肯定相同. const
// index=3 签名编号是:5,4,3,2
//x5 * x4 * x2 / (x5-x3)*(x4-x3)*(x2-x3)
func (v *VerifiableSS) MapShareToNewParams(index int, s []int) share.SPrivKey {
	if len(s) < v.reconstructLimit() {
		panic("reconstructLimit")
	}
	var points []share.SPrivKey
	for i := 1; i <= v.Parameters.ShareCount; i++ {
		points = append(points, share.BigInt2PrivateKey(big.NewInt(int64(i))))
	}
	xi := points[index]
	num := share.BigInt2PrivateKey(big.NewInt(1)) //不需要mod了吧?
	denum := share.BigInt2PrivateKey(big.NewInt(1))
	//num=除了自己的编号，其他的签名人编号连乘
	for i := 0; i < len(s); i++ {
		if s[i] != index {
			num = share.ModMul(num, points[s[i]])
		}
	}
	//除了自己的编号，（其他的签名人编号-我的编号）连乘
	for i := 0; i < len(s); i++ {
		if s[i] != index {
			xj := points[s[i]].Clone()
			share.ModSub(xj, xi)
			share.ModMul(denum, xj)
		}
	}
	//log.Trace(fmt.Sprintf("num=%s,denum=%s", num, denum))
	denum = share.InvertN(denum)
	share.ModMul(num, denum)
	return num
}

/*
   // Performs a Lagrange interpolation in field Zp at the origin
   // for a polynomial defined by `points` and `values`.
   // `points` and `values` are expected to be two arrays of the same size, containing
   // respectively the evaluation points (x) and the value of the polynomial at those point (p(x)).

   // The result is the value of the polynomial at x=0. It is also its zero-degree coefficient.

   // This is obviously less general than `newton_interpolation_general` as we
   // only get a single value, but it is much faster.
*/
func lagrangeInterpolationAtZero(points []share.SPrivKey, values []share.SPrivKey) share.SPrivKey {
	var lagCoef []share.SPrivKey
	for i := 0; i < len(values); i++ {
		xi := points[i]
		yi := values[i]
		num := share.BigInt2PrivateKey(big.NewInt(1))
		denum := share.BigInt2PrivateKey(big.NewInt(1))
		for j := 0; j < len(values); j++ {
			if i != j {
				num = share.ModMul(num, points[j])
			}
		}
		for j := 0; j < len(values); j++ {
			if i != j {
				//denum*=points[j]-xi
				xj := points[j].Clone()
				share.ModSub(xj, xi)
				share.ModMul(denum, xj)
			}
		}
		denum = share.InvertN(denum)
		share.ModMul(num, denum)
		share.ModMul(num, yi) //num*denum*yi
		lagCoef = append(lagCoef, num)
	}
	var result = lagCoef[0].Clone()
	for i := 1; i < len(values); i++ {
		share.ModAdd(result, lagCoef[i])
	}
	return result
}
func (v *VerifiableSS) reconstructLimit() int {
	return v.Parameters.Threshold + 1
}

// Reconstruct : 根据部分信息,还原私钥 const func
func (v *VerifiableSS) Reconstruct(indices []int, shares []share.SPrivKey) share.SPrivKey {
	if len(shares) != len(indices) {
		panic("arg error")
	}
	if len(shares) < v.reconstructLimit() {
		panic("arg error")
	}
	var points []share.SPrivKey
	for _, i := range indices {
		b := big.NewInt(int64(i + 1))
		points = append(points, share.BigInt2PrivateKey(b))
	}
	return lagrangeInterpolationAtZero(points, shares)
}
