package dkg

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

const MAX_SIZE = ^uint16(0)

type Session struct {
	SelfId    uint32
	PartyIds  []uint32
	Size      uint32 // N
	Threshold uint32 // T

	Secret     *secp256k1.Scalar
	Polynomial *commitment.Polynomial
	CommSum    *commitment.CommitmentVector            // componentwise sum of all commitments
	Comms      map[uint32]*commitment.CommitmentVector // commitment vectors of all participants
}

func CreateSession(selfId uint32, partyIds []uint32, size uint32, threshold uint32) *Session {
	if threshold == 0 || size > uint32(MAX_SIZE) || selfId == 0 {
		panic("invalid session parameters")
	}

	return &Session{
		SelfId:    selfId,
		PartyIds:  partyIds,
		Size:      size,
		Threshold: threshold,
		Comms:     make(map[uint32]*commitment.CommitmentVector, size),
	}
}
