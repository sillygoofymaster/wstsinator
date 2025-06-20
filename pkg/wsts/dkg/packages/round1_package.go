package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/pok"
)

type Round1Package struct {
	BasePackage *BasePackage
	PoK         *pok.PoK
	CommVect    *commitment.CommitmentVector
}

func (pkg *Round1Package) GetBase() *BasePackage {
	return pkg.BasePackage
}
