package test

import (
	"fmt"
	"testing"

	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
	"github.com/sillygoofymaster/wstsinator/pkg/wsts/dkg"
	"github.com/sillygoofymaster/wstsinator/pkg/wsts/dkg/packages"
)

func collectPackages(self uint32, pkgs []packages.Packable) []packages.Packable {
	var result []packages.Packable
	for i, pkg := range pkgs {
		base := pkg.GetBase()
		if base == nil || base.From == self || (base.To != nil && base.To.PartyId != self) {
			continue
		}

		// to initiate investigation
		round2pkg, ok := pkg.(*packages.Round2Package)
		if ok && base.From == 2 && base.To.PartyId == 3 && base.To.KeyId == 6 {
			round2pkg.Share = secp256k1.GetRandomScalar()
			pkgs[i] = round2pkg
		}

		result = append(result, pkgs[i])
	}
	return result
}

func mockTransport(sessions map[uint32]*dkg.Session, partyIds []dkg.PartyId) {
	type RoundState struct {
		Round dkg.Roundable
		Pkg   packages.Packable
	}

	n := len(sessions)
	rounds := make(map[uint32]RoundState, n)

	for _, i := range partyIds {
		rounds[i.OwnId] = RoundState{
			Round: &dkg.Round1{Session: sessions[i.OwnId]},
		}
	}

	for {
		var pkgs []packages.Packable

		for _, i := range partyIds {
			pkg := rounds[i.OwnId].Round.Generate()
			pkgs = append(pkgs, pkg...)
		}

		for _, i := range partyIds {
			recvPkgs := collectPackages(uint32(i.OwnId), pkgs)
			pkg, err := rounds[i.OwnId].Round.ProcessAndVerify(recvPkgs)
			if err != nil {
				panic(err)
			}
			if pkg != nil {
				fmt.Printf("participant with id %d finished dkg successfully\n", i)
			}
		}

		done := true

		nextRounds := make(map[uint32]RoundState, n)
		for _, i := range partyIds {
			next := rounds[i.OwnId].Round.NextRound()
			nextRounds[i.OwnId] = RoundState{Round: next}
			if next != nil {
				done = false
			}
		}

		rounds = nextRounds

		if done {
			break
		}
	}
}

func TestWSTS(t *testing.T) {
	threshold := uint32(7)
	partyIds := []dkg.PartyId{
		{
			OwnId:  1,
			KeyIds: []uint32{1, 2, 3},
		},
		{
			OwnId:  2,
			KeyIds: []uint32{4, 5},
		},
		{
			OwnId:  3,
			KeyIds: []uint32{6, 7, 8, 9},
		},
	}

	sessions := make(map[uint32]*dkg.Session, len(partyIds))
	for _, pid := range partyIds {
		sessions[pid.OwnId] = dkg.CreateSession(pid, partyIds, threshold)
	}

	mockTransport(sessions, partyIds)
}
