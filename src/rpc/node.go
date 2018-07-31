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
	reply.Blocks = n.Node.GetBlocks()

	// Broadcast to all peers presence of new peer
	// MINEUR-08
	return nil
}

func (n *NodeRPC) DeliverMessage(args *MessageRPC, reply *int) error {
	// Upon reception of message by a client
	// MINEUR-03
	// CLIENT-07
	// To force delivery in right order, make this a synchronous call
	return n.Node.ReceiveMessage(args.Message, args.Time, args.PeerID, args.MessageType)
}

func (n *NodeRPC) DeliverBlock(args *BlockRPC, reply *int) error {
	// Upon reception of block from a miner
	// To implement
	// MINEUR-09
	// MINEUR-11
	go n.Node.ReceiveBlock(args.Block, args.PeerID)
	return nil
	// If block is valid, halt current work to find block and start a new one including this new one
}

func (n *NodeRPC) GetBlocks(args *GetBlocksRPC, reply *BlocksRPC) error {
	// CLIENT-10
	// To implement
	reply.Blocks = n.Node.GetBlocks()
	return nil
}

func (n *NodeRPC) Disconnect(args *string, reply *int) error {
	go n.Node.CloseConnection(*args)
	return nil
}

func (n *NodeRPC) Connect(args *PeerRPC, reply *int) error {
	n.Node.OpenConnection(args.Host, args.Port)
	return nil
}
