package app

import (
<<<<<<< HEAD
	brpc "LOG735-PG/src/rpc"
	"LOG735-PG/src/node"
	"strings"
	"net"
	"net/http"
	"net/rpc"
	"fmt"
	"log"
=======
	"time"
	"hash"
>>>>>>> first draft creating miner block
)

type Miner struct {
	ID string // i.e. Run-time port associated to container
	blocks []node.Block // MINEUR-07
	peers []string // Slice of IDs
	rpcHandler *brpc.NodeRPC
}

func NewMiner(port, peers string) node.Node {
	m := &Miner{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		strings.Split(peers, " "),
		new(brpc.NodeRPC),
	}
	m.rpcHandler.Node = m
	return m
}
// MINEUR-12

func (m *Miner) SetupRPC(port string) error {
	rpc.Register(m.rpcHandler)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	log.Printf("Listening on TCP port %s\n", port)
	go http.Serve(l, nil)
	return nil
}

func (m Miner) Peer() error {
	for _, peer := range m.peers {
		client, err := brpc.ConnectTo(peer)
		if err != nil {
			return err
		}
		args := &brpc.ConnectionRPC{m.ID}
		var reply brpc.BlocksRPC
		err = client.Call("NodeRPC.Peer", args, &reply)
		if err != nil {
			return err
		}
		if reply.Blocks != nil {
			return fmt.Errorf("Blocks are not defined")
		}
		log.Printf("Successfully peered with node-%s\n", peer)
	}

	return nil
}

func (m Miner) GetBlocks() []node.Block {
	return m.blocks
}

func (m Miner) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// MINEUR-04
	// To implement
	return nil
}

func (m *Miner) CreateBlock() error {
	// MINEUR-10
	// MINEUR-14
	// To implement
<<<<<<< HEAD
	header := &node.Header{}
	err := m.findNounce(header, uint64(0))
=======
	var lastBlockHash hash.Hash64
	lastBlockHash = nil
	if len(m.Chain) > 0 { 
		lastBlockHash = m.Chain[len(m.Chain)-1].Header.Hash
	}
	 
	header := Header{PreviousBlock: lastBlockHash, Date: time.Now()}
	newBlock := Block{Header: header}
	m.Chain = append(m.Chain, newBlock)

	err := m.findNounce(&header, uint64(0))
>>>>>>> first draft creating miner block
	if err != nil {
		return err
	}

	// Broadcast to all peers
	// MINEUR-06
	return nil
}

func (m Miner) findNounce(header *node.Header, difficulty uint64) error {
	// MINEUR-05
	return nil
}