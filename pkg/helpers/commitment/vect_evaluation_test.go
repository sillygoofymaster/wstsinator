package commitment

import (
	"math/rand"
	"testing"

	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
	"github.com/stretchr/testify/assert"
)

func TestVectEvaluate(t *testing.T) {
	max := uint32(^uint16(0))
	threshold := rand.Uint32() % max
	polynomial := GetRandomCoefficients(threshold)
	vect := CreateCommitmentVector(polynomial)
	scalar := secp256k1.GetRandomScalar()
	vecteval1 := vect.Evaluate(scalar)
	vectval2 := vect.EvaluateHorner(scalar)
	check := vecteval1.Equals(vectval2)
	assert.True(t, check == 1)
}
