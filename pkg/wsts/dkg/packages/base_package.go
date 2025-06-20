package packages

type To struct {
	PartyId uint32
	KeyId   uint32
}

type BasePackage struct {
	From uint32
	To   *To
}
