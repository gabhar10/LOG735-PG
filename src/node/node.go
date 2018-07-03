package node

type Node interface {
	SetupRPC(string) error
	Peer() error
	GetBlocks() []Block
}