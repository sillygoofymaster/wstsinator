package pok

import (
	"math/rand"
	"testing"

	"github.com/sillygoofymaster/wstsinator/pkg/internal/secp256k1"
	"github.com/stretchr/testify/assert"
)

func TestPoK(t *testing.T) {
	a0 := secp256k1.GetRandomScalar()

	public := secp256k1.ScalarBaseMultiplication(a0)

	self := rand.Uint32()

	pok := CreatePoK(self, a0)

	assert.True(t, pok.Verify(self, public))
}
