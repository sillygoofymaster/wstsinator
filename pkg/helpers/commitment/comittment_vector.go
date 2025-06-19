package commitment

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type CommitmentVector struct {
	Components []*secp256k1.AffinePoint
}

func CreateCommitmentVector(polynomial *Polynomial) *CommitmentVector {
	n := len(polynomial.Coefficients)
	comps := make([]*secp256k1.AffinePoint, n)

	for i := 0; i < n; i++ {
		comps[i] = secp256k1.ScalarBaseMultiplication(polynomial.Coefficients[i])
	}
	return &CommitmentVector{Components: comps}
}

func AddTwoVectors(A *CommitmentVector, B *CommitmentVector) *CommitmentVector {
	result := make([]*secp256k1.AffinePoint, len(A.Components))
	for i := range A.Components {
		result[i] = secp256k1.AddTwoPoints(A.Components[i], B.Components[i])
	}

	return &CommitmentVector{
		Components: result,
	}
}

func Copy(A *CommitmentVector) *CommitmentVector {
	n := len(A.Components)
	newComps := make([]*secp256k1.AffinePoint, n)
	for i := 0; i < n; i++ {
		temp := new(secp256k1.AffinePoint)
		if A.Components[i] != nil {
			*temp = *A.Components[i]
			newComps[i] = temp
		}
	}
	return &CommitmentVector{
		Components: newComps,
	}
}

func (vect *CommitmentVector) Evaluate(scalar *secp256k1.Scalar) *secp256k1.AffinePoint {
	x := new(secp256k1.Scalar).SetInt(uint32(1))
	result := secp256k1.ScalarMultiplication(x, vect.Components[0])
	x.Mul(scalar)

	var temp = new(secp256k1.AffinePoint)

	for i := 1; i < len(vect.Components); i++ {
		temp = secp256k1.ScalarMultiplication(x, vect.Components[i])
		result = secp256k1.AddTwoPoints(result, temp)
		x.Mul(scalar)
	}

	return result
}

func (vect *CommitmentVector) EvaluateHorner(scalar *secp256k1.Scalar) *secp256k1.AffinePoint {
	n := len(vect.Components)
	result := vect.Components[n-1]

	var temp = new(secp256k1.AffinePoint)

	for i := n - 2; i >= 0; i-- {
		temp = secp256k1.ScalarMultiplication(scalar, result)
		result = secp256k1.AddTwoPoints(vect.Components[i], temp)
	}

	return result
}
