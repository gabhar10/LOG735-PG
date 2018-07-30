package node

type Node interface {
	SetupRPC() error
	Peer() error
	GetBlocks() []Block
	ReceiveMessage(string, string, string, int) error
	ReceiveBlock(Block) error
	Start()
	Connect(string) error
	Disconnect() error
	CloseConnection(string) error
	OpenConnection(string) error
}
