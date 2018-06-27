package rpc

import (
	"LOG735-PG/src/app"
)

type ClientRPC struct {
	client *app.Client
}

func (c *ClientRPC) Peer(args *ConnectionRPC) error {
	// Peer with Client
	// CLIENT-01
	// To implement
	return nil
}

func (c *ClientRPC) DeliverMessage(args *MessageRPC) error {
	// Upon reception of message by a client
	// CLIENT-07
	// To implement
	return nil
}

func (c *ClientRPC) DeliverBlock(args *BlockRPC) error {
	// Upon reception of block from a miner
	// To implement
	return nil
}

func (c *ClientRPC) Disconnect() error {
	// Send to all peers
	// To implement
	return nil
}