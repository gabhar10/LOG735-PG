package rpc

import (
	"hash"
	"LOG735-PG/src/app"
)

type MinerRPC struct {
	// Miner
}

type GetBlocksRPC struct {
	FirstBlock hash.Hash
	LastBlock hash.Hash
}

type BlocksRPC struct {
	Blocks []app.Block
}

func (m *MinerRPC) Peer(args *ConnectionRPC, resp *BlocksRPC) error {
	// Peer with Miner/Client
	// MINEUR-01
	// To implement

	// MINEUR-02
	blocks := make([]app.Block, app.MinBlocksReturnSize)

	// Mutate blocks ...

	resp = &BlocksRPC{
		Blocks: blocks,
	}

	return nil
}

func (m *MinerRPC) DeliverMessage(args *MessageRPC) error {
	// Upon reception of message by a client
	// MINEUR-03
	// To implement
	return nil
}

func (m *MinerRPC) DeliverBlock(args *BlockRPC) error {
	// Upon reception of block from a miner
	// To implement
	return nil
}

func (m *MinerRPC) GetBlocks(args *GetBlocksRPC, res *BlocksRPC) error {
	// CLIENT-10
	// To implement
	return nil
}