package rpc

import (
	"LOG735-PG/src/node"
	"fmt"
	"hash"
	"log"
	"net/rpc"
	"time"
)

const MessageType = 0
const DisconnectionType = 1
const ConnectionType = 2

type (
	ConnectionRPC struct {
		PeerID string
	}

	MessageRPC struct {
		ConnectionRPC
		Message     string
		Time        time.Time
		MessageType int // Indicates if message is chat message [0],
		// Disconnection of a peer [1] or
		// Connection of peer [2]
	}

	BlocksRPC struct {
		ConnectionRPC
		Blocks []node.Block
	}

	BlockRPC struct {
		ConnectionRPC
		Block node.Block
	}

	GetBlocksRPC struct {
		FirstBlock hash.Hash
		LastBlock  hash.Hash
	}
)

func ConnectTo(peer node.Peer) (*rpc.Client, error) {
	maxTries := 5
	var (
		c   *rpc.Client
		err error
	)
	for i := 0; i < maxTries; i++ {
		time.Sleep(time.Second)
		log.Printf("Dialing %s, try #%d\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port), i+1)
		c, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		if err == nil {
			break
		}
	}
	return c, err
}
