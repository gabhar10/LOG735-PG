package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"time"
)

type Miner struct {
	ID         string       // i.e. Run-time port associated to container
	blocks     []node.Block // MINEUR-07
	peers      []string     // Slice of IDs
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
	var lastBlockHash [sha256.Size]byte

	if len(m.blocks) > 0 {
		lastBlockHash = m.blocks[len(m.blocks)-1].Header.Hash
	}

	header := node.Header{PreviousBlock: lastBlockHash, Date: time.Now()}
	newBlock := node.Block{Header: header}
	m.blocks = append(m.blocks, newBlock)

	err := m.findNounce(&header, 2)
	if err != nil {
		return err
	}

	// Broadcast to all peers
	// MINEUR-06
	return nil
}

func (m Miner) findNounce(header *node.Header, difficulty int) error {
	// MINEUR-05
	var hashedHeader [sha256.Size]byte
	var firstCharacters string
	nounce := uint64(0)

	fmt.Println("Difficulty : ", difficulty)

	for {
		header.Nounce = nounce
		hashedHeader = sha256.Sum256([]byte(fmt.Sprintf("%v", header)))
		firstCharacters = string(hashedHeader[:difficulty])

		if strings.Count(firstCharacters, "0") == difficulty {
			break
		}
		nounce++
	}
	fmt.Println("firstCharacters : ", firstCharacters)
	fmt.Println("Nounce : ", nounce)
	header.Hash = hashedHeader
	return nil
}
