package packages

import (
	"github.com/sillygoofymaster/wstsinator/pkg/internal/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/internal/pok"
)

type Round1Package struct {
	BasePackage *BasePackage
	PoK         *pok.PoK
	CommVect    *commitment.CommitmentVector
}

func (*Round1Package) ShouldBroadcast() bool {
	return true
}

func (pkg *Round1Package) GetBase() *BasePackage {
	return pkg.BasePackage
}
