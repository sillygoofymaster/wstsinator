package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type OutputPackage struct {
	Base                    *BasePackage
	SecretShare             map[uint32]*secp256k1.Scalar // a point per key
	PublicVerificationShare map[uint32]*secp256k1.AffinePoint
	GroupPublicKey          *secp256k1.AffinePoint
}

func (pkg *OutputPackage) GetBase() *BasePackage {
	return pkg.Base
}

func NewOutputPackage(from uint32, secshare map[uint32]*secp256k1.Scalar, pubshare map[uint32]*secp256k1.AffinePoint, groupPublicKey *secp256k1.AffinePoint) *OutputPackage {
	base := &BasePackage{
		From: from,
	}

	return &OutputPackage{
		Base:                    base,
		SecretShare:             secshare,
		PublicVerificationShare: pubshare,
		GroupPublicKey:          groupPublicKey,
	}
}
