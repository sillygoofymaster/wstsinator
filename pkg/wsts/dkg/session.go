package dkg

import (
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/commitment"
	"github.com/sillygoofymaster/wstsinator/pkg/helpers/secp256k1"
)

type PartyId struct {
	OwnId  uint32
	KeyIds []uint32
}

const MAX_SIZE = ^uint16(0)

type Session struct {
	SelfId PartyId

	PartyIds   []PartyId
	PartySize  uint32 // Np
	KeySetSize uint32 // Nk
	Threshold  uint32 // T

	Secret     map[uint32]*secp256k1.Scalar
	Polynomial *commitment.Polynomial
	CommSum    *commitment.CommitmentVector            // componentwise sum of all commitments
	Comms      map[uint32]*commitment.CommitmentVector // commitment vectors of all participants
}

func CreateSession(selfId PartyId, partyIds []PartyId, threshold uint32) *Session {
	size := uint32(len(partyIds))
	keySetSize := GetKeySetSize(partyIds)
	if threshold == 0 || keySetSize > uint32(MAX_SIZE) || selfId.OwnId == 0 || keySetSize < threshold {
		panic("invalid session parameters")
	}

	Secret := make(map[uint32]*secp256k1.Scalar, len(selfId.KeyIds))
	for _, key := range selfId.KeyIds {
		Secret[key] = new(secp256k1.Scalar)
	}

	return &Session{
		SelfId: selfId,

		PartyIds:   partyIds,
		PartySize:  size,
		KeySetSize: keySetSize,
		Threshold:  threshold,

		Secret: Secret,
		Comms:  make(map[uint32]*commitment.CommitmentVector, size),
	}
}

func GetKeySetSize(partyIds []PartyId) uint32 {
	var result int
	for _, i := range partyIds {
		result = result + len(i.KeyIds)
	}

	return uint32(result)
}
