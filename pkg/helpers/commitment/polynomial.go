package commitment

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type Polynomial struct {
	Coefficients []*secp256k1.Scalar
}

// sample t random values (a_{i, 0}, ..., a_{i, t-1}) to define
// a degree (t − 1) polynomial f_{i}(x) = a_{i, 0} + a_{i, 1}*x + ... +  a_{i, 1}*x^(t-1)
func GetRandomCoefficients(threshold uint32) *Polynomial {
	coeffs := make([]*secp256k1.Scalar, threshold)

	for i := uint32(0); i < threshold; i++ {
		coeffs[i] = secp256k1.GetRandomScalar()
	}

	return &Polynomial{Coefficients: coeffs}
}

// returns a_{i, 0}
func (polynomial *Polynomial) Secret() *secp256k1.Scalar {
	return polynomial.Coefficients[0]
}

func (polynomial *Polynomial) Evaluate(x *secp256k1.Scalar) *secp256k1.Scalar {
	result := new(secp256k1.Scalar)
	for i := len(polynomial.Coefficients) - 1; i >= 0; i-- {
		result = secp256k1.MultAndAdd(result, x, polynomial.Coefficients[i])
	}
	return result
}
