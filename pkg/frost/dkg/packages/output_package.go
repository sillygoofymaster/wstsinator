package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type OutputPackage struct {
	Base                    *BasePackage
	SecretShare             *secp256k1.Scalar
	PublicVerificationShare *secp256k1.AffinePoint
}

func (pkg *OutputPackage) ShouldBroadcast() bool {
	return false
}

func (pkg *OutputPackage) GetBase() *BasePackage {
	return pkg.Base
}

func NewOutputPackage(from uint32, secshare *secp256k1.Scalar, pubshare *secp256k1.AffinePoint) *OutputPackage {
	base := &BasePackage{
		From: from,
	}

	return &OutputPackage{
		Base:                    base,
		SecretShare:             secshare,
		PublicVerificationShare: pubshare,
	}
}
