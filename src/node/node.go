package node

import (
	"time"
)

type Node interface {
	SetupRPC(string) error
	Peer() error
	GetBlocks() []Block
	ReceiveMessage(string, time.Time)
}
