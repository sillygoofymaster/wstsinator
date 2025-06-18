package commitment

import (
	"math/rand"
	"testing"

	"github.com/sillygoofymaster/wstsinator/pkg/internal/secp256k1"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	max := uint32(^uint16(0))
	threshold := rand.Uint32() % max
	polynomial := GetRandomCoefficients(threshold)
	scalar := secp256k1.GetRandomScalar()

	x := new(secp256k1.Scalar).SetInt(uint32(1))
	temp := new(secp256k1.Scalar)
	temp.Mul2(x, polynomial.Coefficients[0])
	x.Mul(scalar)

	for i := uint32(1); i < threshold; i++ {
		temp = secp256k1.MultAndAdd(x, polynomial.Coefficients[i], temp)
		x.Mul(scalar)
	}

	check := polynomial.Evaluate(scalar)

	assert.True(t, check.Equals(temp))
}
