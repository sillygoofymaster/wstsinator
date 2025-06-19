package pok

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

const CTX_STR = "FROST-SECP256K1"

type PoK struct {
	c *secp256k1.Scalar // c_{i}  = H(i || ctx || a_{i, 0}*G || R), R = k*G
	M *secp256k1.Scalar // M = k + a_{i, 0}*c_{i}
}

func CreatePoK(selfId uint32, a0 *secp256k1.Scalar) *PoK {
	k := secp256k1.GetRandomScalar()

	// R:
	R := secp256k1.ScalarBaseMultiplication(k)

	// M:
	public := secp256k1.ScalarBaseMultiplication(a0)

	c := Hash(selfId, CTX_STR, public, R)

	M := secp256k1.MultAndAdd(a0, c, k)

	return &PoK{
		c: c,
		M: M,
	}
}

func Hash(i uint32, ctx string, public *secp256k1.AffinePoint, R *secp256k1.AffinePoint) *secp256k1.Scalar {
	result := new(secp256k1.Scalar)

	size := binary.Size(uint32(0))
	itobyteslice := make([]byte, size)
	binary.LittleEndian.PutUint32(itobyteslice, uint32(i))

	h := sha256.New()
	_, _ = h.Write(itobyteslice)
	_, _ = h.Write([]byte(ctx))

	buf_public_x := make([]byte, 32)
	buf_public_y := make([]byte, 32)
	public.X.FillBytes(buf_public_x)
	public.Y.FillBytes(buf_public_y)

	_, _ = h.Write(buf_public_x)
	_, _ = h.Write(buf_public_y)

	buf_R_x := make([]byte, 32)
	buf_R_y := make([]byte, 32)
	R.X.FillBytes(buf_R_x)
	R.Y.FillBytes(buf_R_y)
	_, _ = h.Write(buf_R_x)
	_, _ = h.Write(buf_R_y)

	byteresult := h.Sum(nil)
	result.SetByteSlice(byteresult)

	return result
}

func (pok *PoK) Verify(from uint32, a0G *secp256k1.AffinePoint) bool {
	// check  H(l || ctx || a_{l, 0}*G || Rl) ==  H(l || ctx || a_{l, 0}*G || Rl_prime)
	c := pok.c
	Ml := pok.M

	copy_c := new(secp256k1.Scalar).Set(c)

	negative_c := copy_c.Negate()
	negative_ca0G := secp256k1.ScalarMultiplication(negative_c, a0G)
	mlG := secp256k1.ScalarBaseMultiplication(Ml)

	// Rl_prime = (M_{l} - a_{l, 0}*c_{l})*G = (k + a_{l, 0}*c_{l} - a_{l, 0}*c_{l})
	Rl_prime := secp256k1.AddTwoPoints(mlG, negative_ca0G)
	c_prime := Hash(from, CTX_STR, a0G, Rl_prime)

	return c_prime.Equals(c)
}
