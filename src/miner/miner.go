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
	"reflect"
	"strings"
	"sync"
	"time"
)

type Miner struct {
	ID                string // i.e. Run-time port associated to container
	blocks            []node.Block
	miningBlock       node.Block        // MINEUR-07
	peers             []*node.Peer      // Slice of peers
	rpcHandler        *brpc.NodeRPC     // Handler for RPC requests
	IncomingMsgChan   chan node.Message // Channel for incoming messages from other clients
	incomingBlockChan chan node.Block   // Channel for incoming blocks from other miners
	quit              chan bool         // Channel to cancel mining operations
	mutex             *sync.Mutex       // Mutex for synchronization between routines
	waitingList       []node.Message
}

func NewMiner(port string, peers []*node.Peer) node.Node {
	m := &Miner{
		port,
		[]node.Block{},
		node.Block{},
		peers,
		new(brpc.NodeRPC),
		make(chan node.Message, node.MessagesChannelSize),
		make(chan node.Block, node.BlocksChannelSize),
		make(chan bool, 1),
		&sync.Mutex{},
		[]node.Message{},
	}
	m.rpcHandler.Node = m
	return m
}

func (m *Miner) Start() {
	go func() {
		block := m.mining()
		m.blocks = append(m.blocks, block)
		//MINEUR-06
		m.BroadcastBlock(m.blocks)
		m.miningBlock = node.Block{}
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
		client, err := brpc.ConnectTo(*peer)
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

		peer.Conn = client
		log.Printf("Successfully peered with node-%s\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
	}

	return nil
}

func (m *Miner) BroadcastBlock([]node.Block) error {
	if len(m.peers) == 0 {
		return fmt.Errorf("No peers are defined")
	}

	for _, peer := range m.peers {
		if peer.Conn == nil {
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		}
		args := brpc.BlocksRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: m.ID},
			Blocks:        m.blocks}
		var reply *int
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Miner) GetBlocks() []node.Block {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.blocks
}

func (m *Miner) Broadcast(message node.Message) error {
	// DeliverMessage (RPC) to peers
	// To implement

	if len(m.peers) == 0 {
		return fmt.Errorf("No peers are defined")
	}

	for _, peer := range m.peers {
		if peer.Conn == nil {
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		}
		args := brpc.MessageRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: m.ID},
			Message:       message.Content,
			Time:          message.Time}
		var reply *int
		err := peer.Conn.Call("NodeRPC.DeliverMessage", &args, &reply)
		if err != nil {
			return err
		}
	}
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
		msg := <-m.IncomingMsgChan
		m.waitingList = append(m.waitingList, msg)
		messages[i] = msg
	}

	return node.Block{Header: header, Messages: messages}
}

func (m *Miner) ReceiveMessage(content string, temps time.Time, peer string) {
	//MINEUR-04
	found := false
	for _, m := range m.waitingList {
		if reflect.DeepEqual(m, node.Message{peer, content, temps}) {
			found = true
			break
		}
	}
	message := node.Message{peer, content, temps}
	if !found {
		m.Broadcast(message)
	}
	m.IncomingMsgChan <- message
}

func (m *Miner) ReceiveBlock(block node.Block) {
	m.quit <- false
	// MINEUR-05

	valid := false
	//es-ce que le bloc precedant existe dans la chaine
	for i := len(m.blocks) - 1; i >= 0; i-- {
		if block.Header.PreviousBlock == m.blocks[i].Header.Hash {
			valid = true
			break
		}
	}

	//que le hash est correct (bonne difficulter)
	if valid {
		valid = false

		header := node.Header{
			PreviousBlock: block.Header.PreviousBlock,
			Nounce:        block.Header.Nounce,
			Date:          block.Header.Date,
		}
		header.Hash = [sha256.Size]byte{}
		hash := sha256.Sum256([]byte(fmt.Sprintf("%v", header)))
		firstCharacters := string(hash[:node.MiningDifficulty])
		if strings.Count(firstCharacters, "0") == node.MiningDifficulty && hash == block.Header.Hash {
			valid = true
		}
	}
	if !valid {
		m.Start()
	}

	//bloc valide? on doit lajouter a ma chaine et broadcast ?
	//si ils nous manque des messages, on doit aller les cherches sur les autres mineurs ?

	//on enleve les messages du bloc valide qui sont dans le bloc quon etait en train de miner
	if valid {
		for _, receivedMessage := range block.Messages {
			for i, miningMessage := range m.miningBlock.Messages {
				if reflect.DeepEqual(receivedMessage, miningMessage) {
					m.miningBlock.Messages[i] = node.Message{}
					i--
				}
			}
		}
	}

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
		m.miningBlock = m.CreateBlock()
		hashedHeader, _ := m.findingNounce(&m.miningBlock)
		log.Println("Nounce : ", m.miningBlock.Header.Nounce)
		m.miningBlock.Header.Hash = hashedHeader

		//supprimer les messages de la waitingList
		for _, message := range m.miningBlock.Messages {
			for j, unprocMessage := range m.waitingList {
				if reflect.DeepEqual(unprocMessage, message) {
					m.waitingList = append(m.waitingList[:j], m.waitingList[j+1:]...)
					break
				}
			}
		}

		return m.miningBlock
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
