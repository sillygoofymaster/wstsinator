package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type Round2Package struct {
	Base  *BasePackage
	Share *secp256k1.Scalar
}

func (pkg *Round2Package) GetBase() *BasePackage {
	return pkg.Base
}

func CreateRound2Package(From uint32, partyid uint32, keyid uint32, Share *secp256k1.Scalar) *Round2Package {
	To := &To{
		PartyId: partyid,
		KeyId:   keyid,
	}
	base := &BasePackage{
		From: From,
		To:   To,
	}

	return &Round2Package{
		Base:  base,
		Share: Share,
	}
}
