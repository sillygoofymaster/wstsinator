package secp256k1

import (
	"testing"
	//"github.com/sillygoofymaster/wstsinator/pkg/internal/scalar"
)

func TestAdd(t *testing.T) {
	point := new(AffinePoint)

	scalar := GetRandomScalar()
	randompoint := ScalarBaseMultiplication(scalar)

	add := AddTwoPoints(point, randompoint)

	AddTwoPoints(point, add)
}
