package dkg

import (
	"fmt"

	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg/packages"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/pok"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type Round1 struct {
	Session *Session
}

func (round *Round1) Generate() []packages.Packable {
	// every participant Pi samples t random values a_i0,...,a_i1 and uses these values as coefficients to define a degree t−1 polynomial
	round.Session.Polynomial = commitment.GetRandomCoefficients(round.Session.Threshold)

	round.Session.CommSum = commitment.CreateCommitmentVector(round.Session.Polynomial)

	// every Pi computes a proof of knowledge to the corresponding secret a_i0
	a0 := round.Session.Polynomial.Secret()
	round.Session.Secret = new(secp256k1.Scalar).Set(a0)

	pok := pok.CreatePoK(round.Session.SelfId, a0)

	comvect := commitment.Copy(round.Session.CommSum)

	round.Session.Comms[round.Session.SelfId] = commitment.Copy(comvect)

	base := packages.BasePackage{
		From: round.Session.SelfId,
	}

	output := packages.Round1Package{
		BasePackage: &base,
		PoK:         pok,
		CommVect:    comvect,
	}

	return []packages.Packable{&output}
}

// every participant Pi verifies PoK received from every other participant
func (round *Round1) ProcessAndVerify(pkgs []packages.Packable) (packages.Packable, error) {

	if len(pkgs) != int(round.Session.Size)-1 {
		return nil, fmt.Errorf("wrong package amount passed")
	}

	for _, pkg := range pkgs {

		round1Pkg, ok := pkg.(*packages.Round1Package)
		if !ok {
			return nil, fmt.Errorf("wrong package type passed")
		}

		base := round1Pkg.GetBase()
		if base.From == round.Session.SelfId {
			return nil, fmt.Errorf("package addressed to oneself")
		}

		check := round1Pkg.PoK.Verify(base.From, round1Pkg.CommVect.Components[0])
		if !check {
			return nil, fmt.Errorf("fakest knowledge ever proved: %d blamed", base.From)
		}

		round.Session.Comms[base.From] = round1Pkg.CommVect
		round.Session.CommSum = commitment.AddTwoVectors(round.Session.CommSum, round1Pkg.CommVect)
	}

	return nil, nil
}

func (round *Round1) NextRound() Roundable {
	return &Round2{Session: round.Session}
}
