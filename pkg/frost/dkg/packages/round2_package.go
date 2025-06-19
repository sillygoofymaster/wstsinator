package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type Round2Package struct {
	Base  *BasePackage
	Share *secp256k1.Scalar
}

func (*Round2Package) ShouldBroadcast() bool {
	return false
}

func (pkg *Round2Package) GetBase() *BasePackage {
	return pkg.Base
}

func CreateRound2PrepPackage(From uint32, To uint32, Share *secp256k1.Scalar) *Round2Package {
	base := &BasePackage{
		From: From,
		To:   To,
	}

	return &Round2Package{
		Base:  base,
		Share: Share,
	}
}
