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
	ID                string            // i.e. Run-time port associated to container
	blocks            []node.Block      // MINEUR-07
	Peers             []*node.Peer      // Slice of peers
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
		make([]node.Block, node.MinBlocksReturnSize),
		peers,
		new(brpc.NodeRPC),
		make(chan node.Message, node.MessagesChannelSize),
		make(chan node.Block, node.BlocksChannelSize),
		make(chan bool),
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
	}()
}

func (m *Miner) Connect(anchorPort string) error{
	return nil
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
	for _, peer := range m.peers {
		client, err := brpc.ConnectTo(*peer)
		if err != nil {
			return err
		}
		args := &brpc.BlocksRPC{ConnectionRPC: brpc.ConnectionRPC{PeerID: m.ID}, Blocks: m.blocks}
		var reply *int
		err = client.Call("NodeRPC.DeliverBlock", args, &reply)
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


func (m *Miner) ReceiveMessage(content string, temps time.Time, peer string, messageType int) {
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

		//supprimer les messages de la waitingList
		for _, message := range block.Messages {
			for j, unprocMessage := range m.waitingList {
				if reflect.DeepEqual(unprocMessage, message) {
					m.waitingList = append(m.waitingList[:j], m.waitingList[j+1:]...)
					break
				}
			}
		}

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

func (m Miner) Disconnect() error{
	return nil
}

// Close connection of requesting peer
func (m *Miner) CloseConnection(disconnectingPeer string) error{
	for i := 0; i < len(m.peers); i++{
		if m.peers[i].Port == disconnectingPeer{
			log.Printf("Closing connection with %s", disconnectingPeer)
			m.peers[i].Conn.Close()
			m.peers[i] = m.peers[len(m.peers)-1]
			m.peers = m.peers[:len(m.peers)-1]
			break
		}else{
			disconnectionNotice := brpc.MessageRPC{
				brpc.ConnectionRPC{m.ID},
				disconnectingPeer,
				time.Now(),
				brpc.DisconnectionType}
			var reply int
			m.peers[i].Conn.Call("NodeRPC.DeliverMessage", disconnectionNotice, &reply)
		}
	}
	return nil
}

// Open connection to requesting peer (Usually for 2-way communication
func (m *Miner) OpenConnection(connectingPort string) error{
	anchorPeer := &node.Peer{
		Host: fmt.Sprintf("node-%s", connectingPort),
		Port: connectingPort}

	client, err := brpc.ConnectTo(*anchorPeer)
	if err != nil {
		log.Printf("Error while connecting to requesting peer %s", connectingPort)
		return err
	}
	log.Printf("Successfully peered with node-%s\n", connectingPort)
	anchorPeer.Conn = client
	m.peers = append(m.peers, anchorPeer)

	return nil
}
