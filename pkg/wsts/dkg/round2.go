package dkg

import (
	"fmt"
	//"slices"

	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
	"github.com/sillygoofymaster/wstsinator/pkg/wsts/dkg/packages"
)

type Round2 struct {
	Session *Session
}

// each Pi generates a secret share (i, (j, fi(j))) for every key from Nk
// and keeps (i, (l, fl(l))) for every l from own key subset for themselves
func (round *Round2) Generate() []packages.Packable {
	var result []packages.Packable

	for _, party := range round.Session.PartyIds {
		for _, key := range party.KeyIds {
			share := round.Session.Polynomial.Evaluate(secp256k1.IdToScalar(key))
			if round.Session.SelfId.OwnId == party.OwnId {
				round.Session.Secret[key].Add(share)
				continue
			}

			pkg := packages.CreateRound2Package(round.Session.SelfId.OwnId, party.OwnId, key, share)
			result = append(result, pkg)
		}
	}

	return result
}

func (round *Round2) ProcessAndVerify(pkgs []packages.Packable) (packages.Packable, error) {
	n := len(pkgs)
	if n != int(round.Session.PartySize-1)*len(round.Session.SelfId.KeyIds) {
		return nil, fmt.Errorf("wrong package amount passed")
	}

	sharemp := make(map[uint32]map[uint32]*secp256k1.Scalar, n)

	secshare := round.Session.Secret

	for _, pkg := range pkgs {

		round2prepPkg, ok := pkg.(*packages.Round2Package)

		if !ok {
			return nil, fmt.Errorf("wrong package type passed")
		}

		if round2prepPkg.Base.To == nil || round2prepPkg.Base.To.PartyId != round.Session.SelfId.OwnId || round2prepPkg.Base.From == round.Session.SelfId.OwnId {
			return nil, fmt.Errorf("wrong addressee")
		}

		if _, ok := sharemp[round2prepPkg.Base.To.KeyId]; !ok {
			sharemp[round2prepPkg.Base.To.KeyId] = make(map[uint32]*secp256k1.Scalar)
		}
		sharemp[round2prepPkg.Base.To.KeyId][round2prepPkg.Base.From] = round2prepPkg.Share

		secshare[pkg.GetBase().To.KeyId].Add(round2prepPkg.Share)
	}

	pubshare := make(map[uint32]*secp256k1.AffinePoint, len(round.Session.SelfId.KeyIds))
	secshareG := make(map[uint32]*secp256k1.AffinePoint, len(round.Session.SelfId.KeyIds))
	for _, key := range round.Session.SelfId.KeyIds {
		pubshare[key] = round.Session.CommSum.EvaluateHorner(secp256k1.IdToScalar(key))
		secshareG[key] = secp256k1.ScalarBaseMultiplication(secshare[key])
		if pubshare[key].Equals(secshareG[key]) != 1 {
			err := Investigate(sharemp[key], round.Session.Comms, round.Session.PartyIds, round.Session.SelfId.OwnId, key)
			if err != nil {
				return nil, fmt.Errorf("dkg failed: %s", err)
			}
			panic("dkg failed and the faulty party was not identified")
		}
	}
	groupPublicKey := round.Session.CommSum.Components[0]

	output := packages.NewOutputPackage(round.Session.SelfId.OwnId, secshare, secshareG, groupPublicKey)
	return output, nil
}

func (round *Round2) NextRound() Roundable {
	return nil
}

func Investigate(partialSecshares map[uint32]*secp256k1.Scalar, Comms map[uint32]*commitment.CommitmentVector, partyIds []PartyId, self uint32, key uint32) error {
	for _, id := range partyIds {
		if id.OwnId == self {
			continue
		}
		comm := Comms[id.OwnId].EvaluateHorner(secp256k1.IdToScalar(key))
		share := secp256k1.ScalarBaseMultiplication(partialSecshares[id.OwnId])
		check := comm.Equals(share)
		if check != 1 {
			return fmt.Errorf("%d blamed", id)
		}
	}
	return nil
}
