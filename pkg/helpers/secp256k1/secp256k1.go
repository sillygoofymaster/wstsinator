// thin dcred_secp256k1 wrapper

package secp256k1

import (
	"crypto/subtle"
	"math/big"

	dcred_secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

var curve = dcred_secp256k1.S256()

type Scalar = dcred_secp256k1.ModNScalar

type AffinePoint struct {
	X, Y *big.Int
}

func ScalarBaseMultiplication(scalar *Scalar) *AffinePoint {
	arr := scalar.Bytes()
	kBytes := arr[:]

	x, y := curve.ScalarBaseMult(kBytes)
	return &AffinePoint{X: x, Y: y}
}

func ScalarMultiplication(scalar *Scalar, point *AffinePoint) *AffinePoint {
	arr := scalar.Bytes()
	kBytes := arr[:]

	x, y := curve.ScalarMult(point.X, point.Y, kBytes)
	return &AffinePoint{X: x, Y: y}
}

// a*b + c
func MultAndAdd(a *Scalar, b *Scalar, c *Scalar) *Scalar {
	result := new(Scalar)
	result.Mul2(a, b).Add(c)
	return result
}

func AddTwoPoints(A *AffinePoint, B *AffinePoint) *AffinePoint {
	X, Y := curve.Add(A.X, A.Y, B.X, B.Y)
	return &AffinePoint{
		X: X,
		Y: Y,
	}
}

func (point *AffinePoint) Bytes() []byte {
	x := make([]byte, 32)
	y := make([]byte, 32)
	x = point.X.FillBytes(x)
	y = point.Y.FillBytes(y)
	return append(x, y...)
}

func (v *AffinePoint) Equals(u *AffinePoint) int {
	sa, sv := u.Bytes(), v.Bytes()
	return subtle.ConstantTimeCompare(sa, sv)
}
