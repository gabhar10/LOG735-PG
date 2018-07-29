package miner

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
	"sync"
	"time"
)

type Miner struct {
	ID                string            // i.e. Run-time port associated to container
	blocks            []node.Block      // MINEUR-07
	peers             []node.Peer       // Slice of peers
	rpcHandler        *brpc.NodeRPC     // Handler for RPC requests
	IncomingMsgChan   chan node.Message // Channel for incoming messages from other clients
	incomingBlockChan chan node.Block   // Channel for incoming blocks from other miners
	quit              chan bool         // Channel to cancel mining operations
	mutex             *sync.Mutex       // Mutex for synchronization between routines
}

func NewMiner(port string, peers []node.Peer) node.Node {
	m := &Miner{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		peers,
		new(brpc.NodeRPC),
		make(chan node.Message, node.MessagesChannelSize),
		make(chan node.Block, node.BlocksChannelSize),
		make(chan bool),
		&sync.Mutex{},
	}
	m.rpcHandler.Node = m
	return m
}

func (m *Miner) Start() {
	go func() {
		block := m.mining()
		m.blocks = append(m.blocks, block)
	}()
}

// MINEUR-12

func (m *Miner) SetupRPC(port string) error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	rpc.Register(m.rpcHandler)
	rpc.HandleHTTP()
	log.Printf("Listening on TCP port %s\n", port)
	go http.Serve(l, nil)
	return nil
}

func (m *Miner) Peer() error {
	for _, peer := range m.peers {
		client, err := brpc.ConnectTo(peer)
		if err != nil {
			return err
		}
		args := &brpc.ConnectionRPC{PeerID: m.ID}
		var reply brpc.BlocksRPC
		err = client.Call("NodeRPC.Peer", args, &reply)
		if err != nil {
			return err
		}
		if len(reply.Blocks) < node.MinBlocksReturnSize {
			return fmt.Errorf("Returned size of blocks is below %d", node.MinBlocksReturnSize)
		}
		log.Printf("Successfully peered with node-%s\n", peer)
	}

	return nil
}

func (m *Miner) GetBlocks() []node.Block {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.blocks
}

func (m *Miner) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// MINEUR-04
	// To implement

	return nil
}

func (m *Miner) CreateBlock() node.Block {
	// MINEUR-10
	// MINEUR-14
	// To implement
	var lastBlockHash [sha256.Size]byte

	if len(m.blocks) > 0 {
		lastBlockHash = m.blocks[len(m.blocks)-1].Header.Hash
	}

	header := node.Header{PreviousBlock: lastBlockHash, Date: time.Now()}
	var messages [node.BlockSize]node.Message

	for i := 0; i < node.BlockSize; i++ {
		messages[i] = <-m.IncomingMsgChan
	}

	return node.Block{Header: header, Messages: messages}
}

func (m *Miner) ReceiveMessage(content string, temps time.Time, peer string) {
	m.IncomingMsgChan <- node.Message{peer, content, temps}
}

func (m *Miner) ReceiveBlock(block node.Block) {
	m.quit <- false
	// compare receivedBlock with miningBlock and
	// delete messages from miningBlock that are in the receivedBlock if the receivedBlock is valid
	// start another mining if we have len(messageQueue) > node.BlockSize
}

func (m *Miner) mining() node.Block {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	select {
	case <-m.quit:
		return node.Block{}
	default:
		block := m.CreateBlock()
		hashedHeader, _ := m.findingNounce(&block)
		log.Println("Nounce : ", block.Header.Nounce)
		block.Header.Hash = hashedHeader
		return block
	}
}

func (m *Miner) findingNounce(block *node.Block) ([sha256.Size]byte, error) {
	var firstCharacters string
	var hashedHeader [sha256.Size]byte
	nounce := uint64(0)
findingNounce:
	for {
		select {
		case <-m.quit:
			return [sha256.Size]byte{}, fmt.Errorf("Quit")
		default:
			block.Header.Nounce = nounce
			hashedHeader = sha256.Sum256([]byte(fmt.Sprintf("%v", block.Header)))
			firstCharacters = string(hashedHeader[:node.MiningDifficulty])

			if strings.Count(firstCharacters, "0") == node.MiningDifficulty {
				//add semaphore for race condition between mining routines
				break findingNounce
			}
			nounce++
		}
	}
	return hashedHeader, nil
}
