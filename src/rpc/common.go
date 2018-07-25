package rpc

import (
	"LOG735-PG/src/node"
	"fmt"
	"hash"
	"log"
	"net/rpc"
	"time"
)

type (
	ConnectionRPC struct {
		PeerID string
	}

	MessageRPC struct {
		ConnectionRPC
		Message string
		Time    time.Time
	}

	BlocksRPC struct {
		ConnectionRPC
		Blocks []node.Block
	}

	GetBlocksRPC struct {
		FirstBlock hash.Hash
		LastBlock  hash.Hash
	}
)

func ConnectTo(port string) (*rpc.Client, error) {
	maxTries := 5
	var (
		c   *rpc.Client
		err error
	)
	for i := 0; i < maxTries; i++ {
		time.Sleep(time.Second)
		log.Printf("Dialing %s, try #%d\n", port, i+1)
		c, err = rpc.DialHTTP("tcp", fmt.Sprintf("node-%s:%s", port, port))
		if err == nil {
			break
		}
	}
	return c, err
}