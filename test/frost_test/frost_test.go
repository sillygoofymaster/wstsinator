// mock run with a trusted dealer
package test

import (
	"fmt"
	"testing"

	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg"
	"github.com/sillygoofymaster/wstsinator/pkg/frost/dkg/packages"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

func collectPackages(self uint32, pkgs []packages.Packable) []packages.Packable {
	var result []packages.Packable
	for i, pkg := range pkgs {
		base := pkg.GetBase()
		if base == nil || base.From == self || (base.To != 0 && base.To != self) {
			continue
		}

		// to initiate investigation
		round2pkg, ok := pkg.(*packages.Round2Package)
		if ok && base.From == 5 && self == 4 {
			round2pkg.Share = secp256k1.GetRandomScalar()
			pkgs[i] = round2pkg
		}

		result = append(result, pkgs[i])
	}
	return result
}

func mockTransport(sessions map[uint32]*dkg.Session, partyIds []uint32) {
	type RoundState struct {
		Round dkg.Roundable
		Pkg   packages.Packable
	}

	n := len(sessions)
	rounds := make(map[uint32]RoundState, n)

	for _, i := range partyIds {
		rounds[i] = RoundState{
			Round: &dkg.Round1{Session: sessions[i]},
		}
	}

	for {
		var pkgs []packages.Packable

		for _, i := range partyIds {
			pkg := rounds[i].Round.Generate()
			pkgs = append(pkgs, pkg...)
		}

		for _, i := range partyIds {
			recvPkgs := collectPackages(uint32(i), pkgs)
			pkg, err := rounds[i].Round.ProcessAndVerify(recvPkgs)
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
			next := rounds[i].Round.NextRound()
			nextRounds[i] = RoundState{Round: next}
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

func TestFrost(t *testing.T) {
	size := uint32(5)
	threshold := uint32(3)

	partyIds := make([]uint32, size)
	for i := uint32(0); i < size; i++ {
		partyIds[i] = i + 1
	}

	sessions := make(map[uint32]*dkg.Session, size)
	for _, pid := range partyIds {
		sessions[pid] = dkg.CreateSession(pid, partyIds, threshold)
	}

	mockTransport(sessions, partyIds)
}
