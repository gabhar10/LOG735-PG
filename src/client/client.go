package client

import (
	"LOG735-PG/src/block"
)

const BLOCKCHAIN_MAX_SIZE = 10

type Client struct {
	Peers [][]byte // Slice of ports
	Port []byte
	Chain [BLOCKCHAIN_MAX_SIZE]block.Block
}