package packages

type Packable interface {
	ShouldBroadcast() bool
	GetBase() *BasePackage
}
