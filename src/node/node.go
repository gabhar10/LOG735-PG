package node

import (
	"time"
)

type Node interface {
	SetupRPC(string) error
	Peer() error
	GetBlocks() []Block
	ReceiveMessage(string, time.Time, string)
	ReceiveBlock(Block)
	Start()
	Connect(string) error
	Disconnect() error
	CloseConnection(string) error
	OpenConnection(string) error
}
