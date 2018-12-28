package proofs

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

var one = big.NewInt(1)

// ErrMessageTooLong is returned when attempting to encrypt a message which is
// too large for the size of the public key.
var ErrMessageTooLong = errors.New("paillier: message too long for Paillier public key size")

// GenerateKey generates an Paillier keypair of the given bit size using the
// random source random (for example, crypto/rand.Reader).
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
	// First, begin generation of p in the background.
	var p *big.Int
	var errChan = make(chan error, 1)
	go func() {
		var err error
		p, err = rand.Prime(random, bits/2)
		errChan <- err
	}()

	// Now, find a prime q in the foreground.
	q, err := rand.Prime(random, bits/2)
	if err != nil {
		return nil, err
	}

	// Wait for generation of p to complete successfully.
	if err := <-errChan; err != nil {
		return nil, err
	}

	n := new(big.Int).Mul(p, q)
	pp := new(big.Int).Mul(p, p)
	qq := new(big.Int).Mul(q, q)
	pminusone := new(big.Int).Sub(p, one)
	qminusone := new(big.Int).Sub(q, one)
	phi := new(big.Int).Mul(pminusone, qminusone)
	dn := new(big.Int).ModInverse(n, phi)
	dp, dq := crtDecompose(dn, pminusone, qminusone) //D mod (P-1) (or mod Q-1)
	return &PrivateKey{
		PublicKey: PublicKey{
			N:        n,
			NSquared: new(big.Int).Mul(n, n),
			G:        new(big.Int).Add(n, one), // g = n + 1
		},
		p:         p,
		pp:        pp,
		pminusone: pminusone,
		q:         q,
		qq:        qq,
		qminusone: qminusone,
		pinvq:     new(big.Int).ModInverse(p, q),
		dp:        dp,
		dq:        dq,
		hp:        h(p, pp, n),
		hq:        h(q, qq, n),
		n:         n,
	}, nil

}

//const
func crtDecompose(D, p, q *big.Int) (dp, dq *big.Int) {
	dp = new(big.Int).Mod(D, p)
	dq = new(big.Int).Mod(D, q)
	return
}

//const
func crtRecombine(x1, x2, m1, m2, m1inv *big.Int) *big.Int {
	diff := new(big.Int).Sub(x2, x1)
	diff.Mod(diff, m2)
	if diff.Sign() < 0 {
		diff.Add(diff, m2)
	}
	diff.Mul(diff, m1inv)
	diff.Mod(diff, m2)
	x := new(big.Int).Add(x1, diff.Mul(diff, m1))
	return x
}

// PrivateKey represents a Paillier key.
type PrivateKey struct {
	PublicKey
	p         *big.Int
	pp        *big.Int
	pminusone *big.Int
	q         *big.Int
	qq        *big.Int
	qminusone *big.Int
	pinvq     *big.Int
	dp        *big.Int
	dq        *big.Int
	hp        *big.Int
	hq        *big.Int
	n         *big.Int
}

// NewPrivateKey :
func NewPrivateKey(p, q *big.Int) *PrivateKey {
	n := new(big.Int).Mul(p, q)
	pp := new(big.Int).Mul(p, p)
	qq := new(big.Int).Mul(q, q)
	pminusone := new(big.Int).Sub(p, one)
	qminusone := new(big.Int).Sub(q, one)
	phi := new(big.Int).Mul(pminusone, qminusone)
	dn := new(big.Int).ModInverse(n, phi)
	dp, dq := crtDecompose(dn, pminusone, qminusone) //D mod (P-1) (or mod Q-1)
	return &PrivateKey{
		PublicKey: PublicKey{
			N:        n,
			NSquared: new(big.Int).Mul(n, n),
			G:        new(big.Int).Add(n, one), // g = n + 1
		},
		p:         p,
		pp:        pp,
		pminusone: pminusone,
		q:         q,
		qq:        qq,
		qminusone: qminusone,
		pinvq:     new(big.Int).ModInverse(p, q),
		dp:        dp,
		dq:        dq,
		hp:        h(p, pp, n),
		hq:        h(q, qq, n),
		n:         n,
	}
}

// GetPQ :for save
func (p *PrivateKey) GetPQ() (*big.Int, *big.Int) {
	return p.p, p.q
}

// Address represents the public part of a Paillier key.
type PublicKey struct {
	N        *big.Int // modulus
	G        *big.Int // n+1, since p and q are same length
	NSquared *big.Int
}

// Clone :
func (p *PublicKey) Clone() *PublicKey {
	return &PublicKey{
		N:        new(big.Int).Set(p.N),
		G:        new(big.Int).Set(p.G),
		NSquared: new(big.Int).Set(p.NSquared),
	}
}
func h(p *big.Int, pp *big.Int, n *big.Int) *big.Int {
	gp := new(big.Int).Mod(new(big.Int).Sub(one, n), pp)
	lp := l(gp, p)
	hp := new(big.Int).ModInverse(lp, p)
	return hp
}

func l(u *big.Int, n *big.Int) *big.Int {
	//t.Logf("pk=%s", utils.StringInterface(pk, 5))
	return new(big.Int).Div(new(big.Int).Sub(u, one), n)
}

// Encrypt encrypts a plain text represented as a byte array. The passed plain
// text MUST NOT be larger than the modulus of the passed public key.
func Encrypt(pubKey *PublicKey, plainText []byte) ([]byte, error) {
	c, _, err := EncryptAndNonce(pubKey, plainText)
	return c, err
}

// EncryptAndNonce encrypts a plain text represented as a byte array, and in
// addition, returns the nonce used during encryption. The passed plain text
// MUST NOT be larger than the modulus of the passed public key.
func EncryptAndNonce(pubKey *PublicKey, plainText []byte) ([]byte, *big.Int, error) {
	r, err := rand.Int(rand.Reader, pubKey.N)
	if err != nil {
		return nil, nil, err
	}
	//todo bai fix
	//r = big.NewInt(37)
	c, err := EncryptWithNonce(pubKey, r, plainText)
	if err != nil {
		return nil, nil, err
	}

	return c.Bytes(), r, nil
}

// EncryptWithNonce encrypts a plain text represented as a byte array using the
// provided nonce to perform encryption. The passed plain text MUST NOT be
// larger than the modulus of the passed public key.
func EncryptWithNonce(pubKey *PublicKey, r *big.Int, plainText []byte) (*big.Int, error) {
	m := new(big.Int).SetBytes(plainText)
	if pubKey.N.Cmp(m) < 1 { // N < m
		return nil, ErrMessageTooLong
	}

	// c = g^m * r^n mod n^2 = ((m*n+1) mod n^2) * r^n mod n^2
	n := pubKey.N
	c := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mod(new(big.Int).Add(one, new(big.Int).Mul(m, n)), pubKey.NSquared),
			new(big.Int).Exp(r, n, pubKey.NSquared),
		),
		pubKey.NSquared,
	)

	return c, nil
}

// Decrypt decrypts the passed cipher text.
func Decrypt(privKey *PrivateKey, cipherText []byte) ([]byte, error) {
	c := new(big.Int).SetBytes(cipherText)
	if privKey.NSquared.Cmp(c) < 1 { // c < n^2
		return nil, ErrMessageTooLong
	}

	cp := new(big.Int).Exp(c, privKey.pminusone, privKey.pp)
	lp := l(cp, privKey.p)
	mp := new(big.Int).Mod(new(big.Int).Mul(lp, privKey.hp), privKey.p)
	cq := new(big.Int).Exp(c, privKey.qminusone, privKey.qq)
	lq := l(cq, privKey.q)

	mqq := new(big.Int).Mul(lq, privKey.hq)
	mq := new(big.Int).Mod(mqq, privKey.q)
	m := crt(mp, mq, privKey)

	return m.Bytes(), nil
}

func crt(mp *big.Int, mq *big.Int, privKey *PrivateKey) *big.Int {
	u := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Sub(mq, mp), privKey.pinvq), privKey.q)
	m := new(big.Int).Add(mp, new(big.Int).Mul(u, privKey.p))
	return new(big.Int).Mod(m, privKey.n)
}

// AddCipher homomorphically adds together two cipher texts.
// To do this we multiply the two cipher texts, upon decryption, the resulting
// plain text will be the sum of the corresponding plain texts.
func AddCipher(pubKey *PublicKey, cipher1, cipher2 []byte) []byte {
	x := new(big.Int).SetBytes(cipher1)
	y := new(big.Int).SetBytes(cipher2)

	// x * y mod n^2
	return new(big.Int).Mod(
		new(big.Int).Mul(x, y),
		pubKey.NSquared,
	).Bytes()
}

// Add homomorphically adds a passed constant to the encrypted integer
// (our cipher text). We do this by multiplying the constant with our
// ciphertext. Upon decryption, the resulting plain text will be the sum of
// the plaintext integer and the constant.
func Add(pubKey *PublicKey, cipher, constant []byte) []byte {
	c := new(big.Int).SetBytes(cipher)
	x := new(big.Int).SetBytes(constant)

	// c * g ^ x mod n^2
	return new(big.Int).Mod(
		new(big.Int).Mul(c, new(big.Int).Exp(pubKey.G, x, pubKey.NSquared)),
		pubKey.NSquared,
	).Bytes()
}

// Mul homomorphically multiplies an encrypted integer (cipher text) by a
// constant. We do this by raising our cipher text to the power of the passed
// constant. Upon decryption, the resulting plain text will be the product of
// the plaintext integer and the constant.
func Mul(pubKey *PublicKey, cipher []byte, constant []byte) []byte {
	c := new(big.Int).SetBytes(cipher)
	x := new(big.Int).SetBytes(constant)

	// c ^ x mod n^2
	return new(big.Int).Exp(c, x, pubKey.NSquared).Bytes()
}

// ExtractNroot randomness component of a zero ciphertext.
func ExtractNroot(key *PrivateKey, z *big.Int) *big.Int {
	zp, zq := crtDecompose(z, key.p, key.q)
	rp := new(big.Int).Exp(zp, key.dp, key.p)
	rq := new(big.Int).Exp(zq, key.dq, key.q)
	r := crtRecombine(rp, rq, key.p, key.q, key.pinvq)
	return r
}
