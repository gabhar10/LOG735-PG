package rpc

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/block"
)

type ClientRPC struct {
	client *client.Client
}

type ConnectionRPC struct {
	SenderID []byte
}

type MessageRPC struct {
	ConnectionRPC
	Message string
}

type BlockRPC struct {
	ConnectionRPC
	block.Block
}

func (c *ClientRPC) Connect(args *ConnectionRPC) error {
	// Send Connect function to all pre-assigned peers
	// To implement
	return nil
}


func (c *ClientRPC) HandleMessage(args *MessageRPC) error {
	// Upon reception of message by another client
	// To implement
	return nil
}

func (c *ClientRPC) HandleBlock(args *BlockRPC) error {
	// Upon reception of block from a miner
	// To implement
	return nil
}

func (c *ClientRPC) Disconnect() error {
	// Send to all peers
	// To implement
	return nil
}