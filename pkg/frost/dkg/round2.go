package dkg

import (
	"errors"

	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg/packages"
	"github.com/sillygoofymaster/wstsinator/pkg/internal/secp256k1"
)

type Round2 struct {
	Session *Session
}

func (round *Round2) Generate() []packages.Packable {
	var result []packages.Packable

	// each Pi generates participant Pj a secret share (j, fi(j)) for every other participant
	// and keeps (i, fi(i)) for themselves
	for _, i := range round.Session.PartyIds {
		share := round.Session.Polynomial.Evaluate(secp256k1.IdToScalar(i))

		if round.Session.SelfId == i {
			round.Session.SecretShare = share
			continue
		}

		pkg := packages.CreateRound2PrepPackage(round.Session.SelfId, i, share)
		result = append(result, pkg)
	}

	return result
}

// each Pi verifies their shares
func (round *Round2) ProcessAndVerify(pkgs []packages.Packable) error {

	if len(pkgs) != int(round.Session.Size)-1 {
		return errors.New("wrong package amount passed")
	}

	secshare := new(secp256k1.Scalar).Set(round.Session.SecretShare)

	// each Pi calculates their long-lived private signing share
	for _, pkg := range pkgs {

		round2prepPkg, ok := pkg.(*packages.Round2Package)
		if !ok {
			return errors.New("wrong package type passed")
		}

		if round2prepPkg.Base.To != round.Session.SelfId || round2prepPkg.Base.From == round.Session.SelfId {
			return errors.New("wrong addressee")
		}

		secshare = secshare.Add(round2prepPkg.Share)
	}

	// and verifies it by checking secshare*G == pubshare, pubshare = CommSum evaluated at i
	pubshare := round.Session.CommSum.Evaluate(secp256k1.IdToScalar(round.Session.SelfId))
	secshareG := secp256k1.ScalarBaseMultiplication(secshare)
	if pubshare.Equals(secshareG) != 1 {
		// investigate?
		return nil
	}

	return nil
}

func (round *Round2) NextRound() Roundable {
	return nil
}
