package node

import (
	"time"
)

type Node interface {
	SetupRPC() error
	Peer() error
	GetBlocks() []Block
	ReceiveMessage(string, time.Time, string, int)
	ReceiveBlock(Block) error
	Start()
	Connect(string) error
	Disconnect() error
	CloseConnection(string) error
	OpenConnection(string) error
}
