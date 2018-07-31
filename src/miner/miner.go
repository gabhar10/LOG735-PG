package miner

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Miner struct {
	ID                string // i.e. Run-time port associated to container
	blocks            []node.Block
	Peers             []*node.Peer      // Slice of peers
	rpcHandler        *brpc.NodeRPC     // Handler for RPC requests
	IncomingMsgChan   chan node.Message // Channel for incoming messages from other clients
	incomingBlockChan chan node.Block   // Channel for incoming blocks from other miners
	quit              chan bool         // Channel to cancel mining operations
	mutex             *sync.Mutex       // Mutex for synchronization between routines
	startMiningMutex  *sync.Mutex
	waitingList       []node.Message
	malicious         bool
}

func NewMiner(port string, peers []*node.Peer) node.Node {
	log.Println("Miner::Entering NewMiner()")
	defer log.Println("Miner::Leaving NewMiner()")

	m := &Miner{
		port,
		[]node.Block{},
		peers,
		new(brpc.NodeRPC),
		make(chan node.Message, node.MessagesChannelSize),
		make(chan node.Block, node.BlocksChannelSize),
		make(chan bool, 1),
		&sync.Mutex{},
		&sync.Mutex{},
		[]node.Message{},
		strings.Contains(strings.ToLower(os.Getenv("MALICIOUS")), "true"),
	}
	log.Printf("Malicious: %v", m.malicious)
	m.rpcHandler.Node = m
	return m
}

func (m *Miner) Start() {
	log.Printf("Miner-%s::Entering Start()", m.ID)
	defer log.Printf("Miner-%s::Leaving Start()", m.ID)

	go func() {
		for {
			block, err := m.mining()
			// Check if Quit was called
			if block == (node.Block{}) && err == nil {
				continue
			} else if err != nil {
				log.Printf("Error while mining: %v", err)
				return
			}

			m.blocks = append(m.blocks, block)
			m.BroadcastBlock(block)
			m.clearProcessedMessages(&block)
		}
	}()
}

func (m *Miner) Connect(host string, anchorPort string) error {
	log.Printf("Miner-%s::Entering Connect()", m.ID)
	defer log.Printf("Miner-%s::Leaving Connect()", m.ID)

	return nil
}

// MINEUR-12

func (m *Miner) SetupRPC() error {
	log.Printf("Miner-%s::Entering SetupRPC()", m.ID)
	defer log.Printf("Miner-%s::Leaving SetupRPC()", m.ID)

	s := rpc.NewServer()
	s.Register(m.rpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", m.ID))
	if err != nil {
		log.Printf("Error while acquiring listener: %v", err)
		return err
	}
	go s.Accept(listener)
	return nil
}

func (m *Miner) Peer() error {
	log.Printf("Miner-%s::Entering Peer()", m.ID)
	defer log.Printf("Miner-%s::Leaving Peer()", m.ID)

	for _, peer := range m.Peers {
		client, err := brpc.ConnectTo(*peer)
		if err != nil {
			log.Printf("Error while getting RPC: %v", err)
			return err
		}
		args := &brpc.ConnectionRPC{PeerID: peer.Port}
		var reply brpc.BlocksRPC

		err = client.Call("NodeRPC.Peer", args, &reply)
		if err != nil {
			log.Printf("Error while peering: %v", err)
			return err
		}
		peer.Conn = client
		log.Printf("Successfully peered with node-%s\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
	}

	return nil
}

func (m *Miner) BroadcastBlock(b node.Block) error {
	log.Printf("Miner-%s::Entering BroadcastBlock()", m.ID)
	defer log.Printf("Miner-%s::Leaving BroadcastBlock()", m.ID)

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
			ConnectionRPC: brpc.ConnectionRPC{PeerID: m.ID},
			Block:         b}
		var reply *int

		if peer.Conn == nil {
			log.Println("Error: Peer's connection is nil")
			return fmt.Errorf("Error: Peer's connection is nil")
		}
		log.Printf("Sending block to peer %s", peer.Port)
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		log.Printf("Sent block to peer %s", peer.Port)
		if err != nil {
			log.Printf("Error while delivering block: %v", err)
			return err
		}
	}
	return nil
}

func (m *Miner) GetBlocks() []node.Block {
	log.Printf("Miner-%s::Entering GetBlocks()", m.ID)
	defer log.Printf("Miner-%s::Leaving GetBlocks()", m.ID)

	return m.blocks
}

func (m *Miner) Broadcast(message node.Message) error {
	log.Printf("Miner-%s::Entering Broadcast()", m.ID)
	defer log.Printf("Miner-%s::Leaving Broadcast()", m.ID)
	// DeliverMessage (RPC) to peers
	// To implement

	if len(m.Peers) == 0 {
		log.Println("No peers are defined. Exiting.")
		return nil
	}

	for _, peer := range m.Peers {
		if peer.Conn == nil {
			log.Printf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
		}
		args := brpc.MessageRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: message.Peer},
			Message:       message.Content,
			Time:          message.Time}
		var reply *int

		err := peer.Conn.Call("NodeRPC.DeliverMessage", &args, &reply)
		if err != nil {
			log.Printf("Error while delivering message: %v", err)
			return err
		}
	}
	return nil
}

func (m *Miner) CreateBlock() node.Block {
	log.Printf("Miner-%s::Entering CreateBlock()", m.ID)
	defer log.Printf("Miner-%s::Leaving CreateBlock()", m.ID)
	// MINEUR-10
	// MINEUR-14
	// To implement
	var lastBlockHash [sha256.Size]byte

	if len(m.blocks) > 0 {
		lastBlockHash = m.blocks[len(m.blocks)-1].Header.Hash
	}

	header := node.Header{PreviousBlock: lastBlockHash, Date: time.Now().Format(time.RFC3339)}
	var messages [node.BlockSize]node.Message

	for i := 0; i < node.BlockSize; i++ {
		msg := <-m.IncomingMsgChan
		messages[i] = msg
	}

	return node.Block{Header: header, Messages: messages}
}

func (m *Miner) ReceiveMessage(content, temps, peer string, messageType int) error {
	log.Printf("Miner-%s::Entering ReceiveMessage()", m.ID)
	defer log.Printf("Miner-%s::Leaving ReceiveMessage()", m.ID)

	// If malicious miner, alternate content of received messages
	if m.malicious {
		contentByte := []byte(content)
		for i := 0; i < len(contentByte); i++ {
			contentByte[i]++
		}
		content = string(contentByte)
	}

	// Check if we don't alrady have the message in the waiting list
	msg := node.Message{
		Peer:    peer,
		Content: content,
		Time:    temps,
	}
	for _, m := range m.waitingList {
		if reflect.DeepEqual(m, msg) {
			log.Printf("Message \"%s\" from peer \"%s\" is already in waiting list from %s", content, peer, msg.Time)
			return nil
		}
	}
	log.Printf("Appending message \"%s\" from peer \"%s\" in waiting list at %s", msg.Content, msg.Peer, msg.Time)
	log.Printf("Locking waitingList mutex")
	m.mutex.Lock()
	m.waitingList = append(m.waitingList, msg)
	log.Printf("Unlocking waitingList mutex")
	m.mutex.Unlock()

	err := m.Broadcast(msg)
	if err != nil {
		log.Printf("Error while broadcasting to peer: %v", err)
		return err
	}
	m.IncomingMsgChan <- msg
	return nil
}

func (m *Miner) ReceiveBlock(block node.Block, peer string) error {
	log.Printf("Miner-%s::Entering ReceiveBlock()", m.ID)
	defer log.Printf("Miner-%s::Leaving ReceiveBlock()", m.ID)

	log.Printf("Received block %v from \"%s\"", block, peer)
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

		if block.Header.Hash != hash {
			log.Printf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
			return fmt.Errorf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
		}

		firstCharacters := string(hash[:node.MiningDifficulty])
		if strings.Count(firstCharacters, "0") == node.MiningDifficulty && hash == block.Header.Hash {
			log.Println("Received block's hash is valid! Stopping current mining operation.")
		} else {
			log.Println("Received block's hash is not valid! Discarding block")
			return nil
		}
	} else {
		log.Println("Block does not exist in my chain! Discard")
		return nil
	}

	//on enleve les messages du bloc valide qui sont dans le bloc quon etait en train de miner

	m.clearProcessedMessages(&block)
	m.quit <- true
	// Malicious nodes do not broadcast new valid blocks
	if !m.malicious {
		err := m.BroadcastBlock(block)
		if err != nil {
			log.Printf("Error while broadcasting block: %v", err)
			return err
		}
	}

	log.Println("Appending block to the chain")
	// TODO: Use mutex?
	tempBlocks := append(m.blocks, block)
	m.blocks = tempBlocks
	return nil
}

func (m *Miner) mining() (node.Block, error) {
	log.Printf("Miner-%s::Entering mining()", m.ID)
	defer log.Printf("Miner-%s::Leaving mining()", m.ID)

	// This method should only be called via Start()
	block := m.CreateBlock()
	hashedHeader, err := m.findingNounce(&block)
	if err != nil {
		log.Printf("Error while finding nounce: %v", err)
		return node.Block{}, fmt.Errorf("Error while finding nounce: %v", err)
	}
	// Did quit get invoked?
	if hashedHeader == ([sha256.Size]byte{}) {
		log.Println("Quit was invoked. Leaving mining")
		return node.Block{}, nil
	}

	log.Println("Nounce : ", block.Header.Nounce)
	block.Header.Hash = hashedHeader

	log.Printf("block: %v", block)
	return block, nil
}

func (m *Miner) findingNounce(block *node.Block) ([sha256.Size]byte, error) {
	log.Printf("Miner-%s::Entering findingNounce()", m.ID)
	defer log.Printf("Miner-%s::Leaving findingNounce()", m.ID)

	var firstCharacters string
	var hashedHeader [sha256.Size]byte
	nounce := uint64(0)
findingNounce:
	for {
		select {
		case <-m.quit:
			log.Println("Quitting!")
			return [sha256.Size]byte{}, nil
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
		if m.malicious {
			time.Sleep(time.Nanosecond * 6000)
		}
	}
	log.Printf("Took %d tries to find nounce", nounce)
	return hashedHeader, nil
}

func (m Miner) Disconnect() error {
	log.Printf("Miner-%s::Entering Disconnect()", m.ID)
	defer log.Printf("Miner-%s::Leaving Disconnect()", m.ID)

	return nil
}

// Close connection of requesting peer
func (m *Miner) CloseConnection(disconnectingPeer string) error {
	log.Printf("Miner-%s::Entering CloseConnection()", m.ID)
	defer log.Printf("Miner-%s::Leaving CloseConnection()", m.ID)

	for i := 0; i < len(m.Peers); i++ {
		if m.Peers[i].Port == disconnectingPeer {
			log.Printf("Closing connection with %s", disconnectingPeer)
			if m.Peers[i].Conn != nil {
				m.Peers[i].Conn.Close()
			}
			m.Peers = append(m.Peers[:i], m.Peers[i+1:]...)
			break
		}
	}
	return nil
}

// Open connection to requesting peer (Usually for 2-way communication
func (m *Miner) OpenConnection(host string, connectingPort string) error {
	log.Printf("Miner-%s::Entering OpenConnection()", m.ID)
	defer log.Printf("Miner-%s::Leaving OpenConnection()", m.ID)

	anchorPeer := &node.Peer{
		Host: host,
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
	log.Printf("Miner-%s::Entering clearProcessedMessages()", m.ID)
	defer log.Printf("Miner-%s::Leaving clearProcessedMessages()", m.ID)

	//supprimer les messages de la waitingList
	log.Printf("Locking waitingList mutex")
	m.mutex.Lock()
	for _, message := range block.Messages {
		for j, unprocMessage := range m.waitingList {
			if reflect.DeepEqual(unprocMessage, message) {
				log.Printf("Message \"%s\" from peer \"%s\" is in waiting list. Removing it.", message.Content, message.Peer)
				m.waitingList = append(m.waitingList[:j], m.waitingList[j+1:]...)
				break
			}
		}
	}
	log.Printf("Unlocking waitingList mutex")
	m.mutex.Unlock()
}
