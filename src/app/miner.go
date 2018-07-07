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

type Message struct {
	Content string
	Time    time.Time
}

type Miner struct {
	ID                string       // i.e. Run-time port associated to container
	blocks            []node.Block // MINEUR-07
	peers             []string     // Slice of IDs
	rpcHandler        *brpc.NodeRPC
	incomingMsgChan   chan Message
	incomingBlockChan chan node.Block
	quit              chan bool
}

func NewMiner(port, peers string) node.Node {
	m := &Miner{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		strings.Split(peers, " "),
		new(brpc.NodeRPC),
		make(chan Message, 100),
		make(chan node.Block, 10),
		make(chan bool),
	}
	m.rpcHandler.Node = m
	return m
}

func (m *Miner) Start() {
	go m.mining()
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

func (m *Miner) CreateBlock() node.Block {
	// MINEUR-10
	// MINEUR-14
	// To implement
	var lastBlockHash [sha256.Size]byte

	if len(m.blocks) > 0 {
		lastBlockHash = m.blocks[len(m.blocks)-1].Header.Hash
	}

	header := node.Header{PreviousBlock: lastBlockHash, Date: time.Now()}
	var messages [node.BlockSize]Message

	for i := 0; i < node.BlockSize; i++ {
		messages[i] = <-m.incomingMsgChan
	}

	return node.Block{Header: header}
}

func (m Miner) ReceiveMessage(content string, temps time.Time) {
	//when recieving message, add it to messageQueue
	m.incomingMsgChan <- Message{content, temps}
}

func (m Miner) ReceiveBlock(block node.Block) {
	m.quit <- false
	// compare receivedBlock with miningBlock and
	// delete messages from miningBlock that are in the receivedBlock if the receivedBlock is valid
	// start another mining if we have len(messageQueue) > node.BlockSize
}

func (m *Miner) mining() {
	var hashedHeader [sha256.Size]byte
	var firstCharacters string
	nounce := uint64(0)
	block := m.CreateBlock()

	fmt.Println("Difficulty : ", node.MiningDifficulty)
findingNounce:
	for {
		select {
		case <-m.quit:
			return
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

	fmt.Println("firstCharacters : ", firstCharacters)
	fmt.Println("Nounce : ", nounce)
	block.Header.Hash = hashedHeader
	m.blocks = append(m.blocks, block)
	if len(m.quit) > 0 {
		<-m.quit
	}
}
