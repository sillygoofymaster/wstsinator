package dkg

import (
	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg/packages"
)

// structured as described in https://eprint.iacr.org/2020/852.pdf
type Roundable interface {
	Generate() []packages.Packable
	ProcessAndVerify([]packages.Packable) error
	NextRound() Roundable
}
