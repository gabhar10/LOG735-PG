package miner

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Miner struct {
	ID     string // i.e. Run-time port associated to container
	blocks []node.Block
	//miningBlock       *node.Block       // MINEUR-07
	Peers             []*node.Peer      // Slice of peers
	rpcHandler        *brpc.NodeRPC     // Handler for RPC requests
	IncomingMsgChan   chan node.Message // Channel for incoming messages from other clients
	incomingBlockChan chan node.Block   // Channel for incoming blocks from other miners
	quit              chan bool         // Channel to cancel mining operations
	mutex             *sync.Mutex       // Mutex for synchronization between routines
	waitingList       []node.Message
}

func NewMiner(port string, peers []*node.Peer) node.Node {
	log.Println("Entering NewMiner()")
	defer log.Println("Leaving NewMiner()")

	m := &Miner{
		port,
		[]node.Block{},
		//nil,
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
	log.Println("Entering Start()")
	defer log.Println("Leaving Start()")

	go func() {
		block := m.mining()
		//m.miningBlock = &block

		//MINEUR-06

		m.blocks = append(m.blocks, block)

		m.BroadcastBlock(block)
		m.mutex.Lock()
		log.Println("Locking mutex")
		m.clearProcessedMessages(&block)
		m.mutex.Unlock()
		log.Println("Unlocking mutex")
	}()
}

func (m *Miner) Connect(anchorPort string) error {
	log.Println("Entering Connect()")
	defer log.Println("Leaving Connect()")

	return nil
}

// MINEUR-12

func (m *Miner) SetupRPC() error {
	log.Println("Entering SetupRPC()")
	defer log.Println("Leaving SetupRPC()")

	s := rpc.NewServer()
	s.Register(m.rpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", m.ID))
	if err != nil {
		return err
	}
	go s.Accept(listener)
	return nil
}

func (m *Miner) Peer() error {
	log.Println("Entering Peer()")
	defer log.Println("Leaving Peer()")

	for _, peer := range m.Peers {
		client, err := brpc.ConnectTo(*peer)
		if err != nil {
			return err
		}
		args := &brpc.ConnectionRPC{PeerID: peer.Port}
		var reply brpc.BlocksRPC

		err = client.Call("NodeRPC.Peer", args, &reply)
		if err != nil {
			return err
		}
		peer.Conn = client
		log.Printf("Successfully peered with node-%s\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
	}

	return nil
}

func (m *Miner) BroadcastBlock(b node.Block) error {
	log.Println("Entering BroadcastBlock()")
	defer log.Println("Leaving BroadcastBlock()")

	if len(m.Peers) == 0 {
		log.Println("No peers are defined. Exiting.")
		return nil
	}

	for _, peer := range m.Peers {
		if peer.Conn == nil {
			log.Println("No connection!")
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		}
		args := brpc.BlockRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: peer.Port},
			Block:         b}
		var reply *int
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Miner) GetBlocks() []node.Block {
	log.Println("Entering GetBlocks()")
	defer log.Println("Leaving GetBlocks()")

	return m.blocks
}

func (m *Miner) Broadcast(message node.Message) error {
	log.Println("Entering Broadcast()")
	defer log.Println("Leaving Broadcast()")
	// DeliverMessage (RPC) to peers
	// To implement

	if len(m.Peers) == 0 {
		return fmt.Errorf("No peers are defined")
	}

	for _, peer := range m.Peers {
		if peer.Conn == nil {
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		}
		args := brpc.MessageRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: peer.Port},
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
	log.Println("Entering CreateBlock()")
	defer log.Println("Leaving CreateBlock()")
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
		messages[i] = msg
	}

	return node.Block{Header: header, Messages: messages}
}

func (m *Miner) ReceiveMessage(content string, temps time.Time, peer string, messageType int) {
	log.Println("Entering ReceiveMessage()")
	defer log.Println("Leaving ReceiveMessage()")

	//MINEUR-04

	// Check if we don't alrady have the message in the waiting list
	msg := node.Message{
		Peer:    peer,
		Content: content,
		Time:    temps,
	}
	for _, m := range m.waitingList {
		if reflect.DeepEqual(m, msg) {
			return
		}
	}
	log.Printf("Appending message \"%s\" from peer \"%s\" in waiting list", msg.Content, msg.Peer)
	log.Println("Locking mutex")
	m.mutex.Lock()
	m.waitingList = append(m.waitingList, msg)
	m.mutex.Unlock()
	log.Println("Unlocking mutex")

	m.Broadcast(msg)
	m.IncomingMsgChan <- msg
}

func (m *Miner) ReceiveBlock(block node.Block) error {
	log.Println("Miner::Entering ReceiveBlock()")
	defer log.Println("Miner::Leaving ReceiveBlock()")

	// Do we already have this block in the chain?
	for _, b := range m.blocks {
		if reflect.DeepEqual(b, block) {
			log.Println("We already have this block in the chain. Discarding block")
			return nil
		}
	}

	valid := false
	//es-ce que le bloc precedant existe dans la chaine
	for i := len(m.blocks) - 1; i >= 0; i-- {
		if block.Header.PreviousBlock == m.blocks[i].Header.Hash {
			valid = true
			break
		}
	}

	// We don't have any blocks and the received one is a genesis block
	if block.Header.PreviousBlock == ([sha256.Size]byte{}) && len(m.blocks) == 0 {
		valid = true
	}

	//que le hash est correct (bonne difficulter)
	if valid {
		header := node.Header{
			PreviousBlock: block.Header.PreviousBlock,
			Nounce:        block.Header.Nounce,
			Date:          block.Header.Date,
		}
		header.Hash = [sha256.Size]byte{}
		hash := sha256.Sum256([]byte(fmt.Sprintf("%v", header)))
		firstCharacters := string(hash[:node.MiningDifficulty])
		if strings.Count(firstCharacters, "0") == node.MiningDifficulty && hash == block.Header.Hash {
			log.Println("Received block's hash is valid! Stopping current mining operation.")
			m.quit <- false
			// MINEUR-05
		} else {
			log.Println("Received block's hash is not valid! Discarding block")
			return nil
		}

	} else {
		log.Println("Block does not exist in my chain! Ask my peers for the missing range of blocks")
		// TODO: Implement fetching range of blocks from peers
		return nil
	}

	//TODO: add timestamp gap checking between messages

	//on enleve les messages du bloc valide qui sont dans le bloc quon etait en train de miner

	m.mutex.Lock()
	m.clearProcessedMessages(&block)
	m.mutex.Unlock()
	err := m.BroadcastBlock(block)
	if err != nil {
		return err
	}

	m.Start()

	return nil
}

func (m *Miner) mining() node.Block {
	log.Println("Entering mining()")
	defer log.Println("Leaving mining()")

	select {
	case <-m.quit:
		return node.Block{}
	default:
		block := m.CreateBlock()
		hashedHeader, err := m.findingNounce(&block)
		if err != nil {
			log.Printf("Error while finding nounce: %v", err)
			return node.Block{}
		}
		log.Println("Nounce : ", block.Header.Nounce)
		block.Header.Hash = hashedHeader

		log.Printf("block: %v", block)
		return block
	}
}

func (m *Miner) findingNounce(block *node.Block) ([sha256.Size]byte, error) {
	log.Println("Entering findingNounce()")
	defer log.Println("Leaving findingNounce()")

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

func (m Miner) Disconnect() error {
	log.Println("Entering Disconnect()")
	defer log.Println("Leaving Disconnect()")

	return nil
}

// Close connection of requesting peer
func (m *Miner) CloseConnection(disconnectingPeer string) error {
	log.Println("Entering CloseConnection()")
	defer log.Println("Leaving CloseConnection()")

	for i := 0; i < len(m.Peers); i++ {
		if m.Peers[i].Port == disconnectingPeer {
			log.Printf("Closing connection with %s", disconnectingPeer)
			if m.Peers[i].Conn != nil {
				m.Peers[i].Conn.Close()
			}
			m.Peers = append(m.Peers[:i], m.Peers[i+1:]...)
			break
		} else {
			disconnectionNotice := brpc.MessageRPC{
				brpc.ConnectionRPC{m.ID},
				disconnectingPeer,
				time.Now(),
				brpc.DisconnectionType}
			var reply int
			if m.Peers[i].Conn == nil {
				return fmt.Errorf("Miner does not have a connection with peer %d", i)
			}
			m.Peers[i].Conn.Call("NodeRPC.DeliverMessage", disconnectionNotice, &reply)
		}
	}
	return nil
}

// Open connection to requesting peer (Usually for 2-way communication
func (m *Miner) OpenConnection(connectingPort string) error {
	log.Println("Entering OpenConnection()")
	defer log.Println("Leaving OpenConnection()")

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
	m.Peers = append(m.Peers, anchorPeer)

	return nil
}

func (m *Miner) clearProcessedMessages(block *node.Block) {
	//supprimer les messages de la waitingList
	for _, message := range block.Messages {
		for j, unprocMessage := range m.waitingList {
			if reflect.DeepEqual(unprocMessage, message) {
				log.Printf("Message \"%s\" from peer \"%s\" is in waiting list. Removing it.", message.Content, message.Peer)
				m.waitingList = append(m.waitingList[:j], m.waitingList[j+1:]...)
				break
			}
		}
	}
}
