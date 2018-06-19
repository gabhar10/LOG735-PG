package rpc

import (
	"hash"
	"LOG735-PG/src/app"
)

type MinerRPC struct {
	miner *app.Miner
}

type GetBlocksRPC struct {
	FirstBlock hash.Hash
	LastBlock hash.Hash
}

type BlocksRPC struct {
	Blocks []app.Block
}

func (m *MinerRPC) Peer(args *ConnectionRPC, resp *BlocksRPC) error {
	// MINEUR-13
	// Peer with Miner/Client
	// MINEUR-01
	// To implement

	// MINEUR-02
	blocks := make([]app.Block, app.MinBlocksReturnSize)

	// Mutate blocks ...

	resp = &BlocksRPC{
		Blocks: blocks,
	}

	// Broadcast to all peers presence of new peer
	// MINEUR-08
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
	// MINEUR-09
	// MINEUR-11
	// If block is valid, halt current work to find block and start a new one including this new one
	return nil
}

func (m *MinerRPC) GetBlocks(args *GetBlocksRPC, res *BlocksRPC) error {
	// CLIENT-10
	// To implement
	return nil
}