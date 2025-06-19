package secp256k1

import (
	"crypto/rand"
)

func GetRandomScalar() *Scalar {
	result := new(Scalar)
	tempbytes := make([]byte, 32)
	if _, err := rand.Read(tempbytes); err != nil {
		panic(err)
	}
	result.SetByteSlice(tempbytes)
	return result
}

func IdToScalar(id uint32) *Scalar {
	result := new(Scalar)
	result.SetInt(id)
	return result
}
