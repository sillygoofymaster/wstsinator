package dkg

import (
	"fmt"

	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg/packages"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type Round2 struct {
	Session *Session
}

// each Pi generates participant Pj a secret share (j, fi(j)) for every other participant
// and keeps (i, fi(i)) for themselves
func (round *Round2) Generate() []packages.Packable {
	var result []packages.Packable

	for _, i := range round.Session.PartyIds {
		share := round.Session.Polynomial.Evaluate(secp256k1.IdToScalar(i))
		if round.Session.SelfId == i {
			round.Session.Secret = share
			continue
		}

		pkg := packages.CreateRound2Package(round.Session.SelfId, i, share)
		result = append(result, pkg)
	}

	return result
}

// each Pi calculates their long-lived private signing share
// and verifies it by checking secshare*G == pubshare (pubshare = CommSum evaluated at i)
func (round *Round2) ProcessAndVerify(pkgs []packages.Packable) (packages.Packable, error) {

	n := len(pkgs)
	if n != int(round.Session.Size)-1 {
		return nil, fmt.Errorf("wrong package amount passed")
	}

	sharemp := make(map[uint32]*secp256k1.Scalar, n)

	secshare := new(secp256k1.Scalar).Set(round.Session.Secret)

	for _, pkg := range pkgs {

		round2prepPkg, ok := pkg.(*packages.Round2Package)

		if !ok {
			return nil, fmt.Errorf("wrong package type passed")
		}

		sharemp[round2prepPkg.Base.From] = round2prepPkg.Share

		if round2prepPkg.Base.To != round.Session.SelfId || round2prepPkg.Base.From == round.Session.SelfId {
			return nil, fmt.Errorf("wrong addressee")
		}

		secshare = secshare.Add(round2prepPkg.Share)
	}

	pubshare := round.Session.CommSum.EvaluateHorner(secp256k1.IdToScalar(round.Session.SelfId))
	secshareG := secp256k1.ScalarBaseMultiplication(secshare)
	if pubshare.Equals(secshareG) != 1 {
		err := Investigate(sharemp, round.Session.Comms, round.Session.PartyIds, round.Session.SelfId)
		if err != nil {
			return nil, fmt.Errorf("dkg failed: %s", err)
		}
		panic("dkg failed and the faulty party was not identified")
	}
	groupPublicKey := round.Session.CommSum.Components[0]

	output := packages.NewOutputPackage(round.Session.SelfId, secshare, secshareG, groupPublicKey)
	return output, nil
}

func (round *Round2) NextRound() Roundable {
	return nil
}

func Investigate(partialSecshares map[uint32]*secp256k1.Scalar, Comms map[uint32]*commitment.CommitmentVector, partyIds []uint32, self uint32) error {
	for _, id := range partyIds {
		if id == self {
			continue
		}
		comm := Comms[id].EvaluateHorner(secp256k1.IdToScalar(self))
		share := secp256k1.ScalarBaseMultiplication(partialSecshares[id])
		check := comm.Equals(share)
		if check != 1 {
			return fmt.Errorf("%d blamed", id)
		}
	}
	return nil
}
