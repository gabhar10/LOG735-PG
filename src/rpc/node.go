package rpc

import (
	"LOG735-PG/src/node"
)

type NodeRPC struct {
	Node node.Node
}

func (n *NodeRPC) Peer(args *ConnectionRPC, reply *BlocksRPC) error {
	// Peer with Client
	// CLIENT-01
	// To implement
	// MINEUR-13
	// Peer with Miner/Client
	// MINEUR-01
	// To implement

	// MINEUR-02

	// Mutate blocks ...

	reply = &BlocksRPC{
		Blocks: n.Node.GetBlocks(),
	}

	// Broadcast to all peers presence of new peer
	// MINEUR-08
	return nil
}


func (n *NodeRPC) DeliverMessage(args *MessageRPC, reply *int) error {
	// Upon reception of message by a client
	// MINEUR-03
	// CLIENT-07
	// To implement
	//log.Printf("NODE-RPC : I received %s from %s\n", args.Message, args.PeerID)
	n.Node.ReceiveMessage(args.Message, args.Time, args.PeerID)
	return nil
}

func (n *NodeRPC) DeliverBlock(args *BlocksRPC, reply *int) error {
	// Upon reception of block from a miner
	// To implement
	// MINEUR-09
	// MINEUR-11
	n.Node.ReceiveBlock(args.Blocks[len(args.Blocks)-1])
	// If block is valid, halt current work to find block and start a new one including this new one
	return nil
}

func (n *NodeRPC) GetBlocks(args *GetBlocksRPC, reply *BlocksRPC) error {
	// CLIENT-10
	// To implement
	return nil
}

func (n *NodeRPC) Disconnect(args *string, reply *int) error {
	// Send to all peers
	// To implement
	return nil
}
