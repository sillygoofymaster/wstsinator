package dkg

import (
	"github.com/sillygoofymaster/wstsinator/pkg/wsts/dkg/packages"
)

type Roundable interface {
	Generate() []packages.Packable
	ProcessAndVerify([]packages.Packable) (packages.Packable, error)
	NextRound() Roundable
}
