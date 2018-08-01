package client

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

type (
	Client struct {
		Host              string
		ID                string            // i.e. Run-time port associated to container
		blocks            []node.Block      // Can be a subset of the full chain
		Peers             []*node.Peer      // Slice of Peers
		rpcHandler        *brpc.NodeRPC     // Handler for RPC calls
		uiChannel         chan node.Message // Channel to send message from node to chat application
		nodeChannel       chan node.Message // Channel to send message from chat application to node
		msgLoopChan       chan string
		msgLoopRunning    bool
		blocksMutex       *sync.Mutex
		trafficGeneration bool
	}
)

func NewClient(host string, port string, peers []*node.Peer, uiChannel chan node.Message, nodeChannel chan node.Message,
	trafficGeneration bool) node.Node {
	log.Println("Client::Entering NewClient()")
	defer log.Println("Client::Leaving NewClient()")

	c := &Client{
		host,
		port,
		[]node.Block{},
		peers,
		new(brpc.NodeRPC),
		uiChannel,
		nodeChannel,
		make(chan string),
		false,
		&sync.Mutex{},
		trafficGeneration,
	}
	c.rpcHandler.Node = c
	return c
}

func (c *Client) Start() {
	log.Printf("Client-%s::Entering Start()", c.ID)
	defer log.Printf("Client-%s::Leaving Start()", c.ID)

	go c.StartMessageLoop()
}

func (c *Client) GetChain() {
	for _, block := range c.blocks {
		c.ParseBlock(block)
	}
}

func (c *Client) SetupRPC() error {
	log.Printf("Client-%s::Entering SetupRPC()", c.ID)
	defer log.Printf("Client-%s::Leaving SetupRPC()", c.ID)

	s := rpc.NewServer()
	s.Register(c.rpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", c.ID))
	if err != nil {
		log.Printf("Error while acquiring listener: %v", err)
		return err
	}
	go s.Accept(listener)
	return nil
}

func (c *Client) Peer() error {
	log.Printf("Client-%s::Entering Peer()", c.ID)
	defer log.Printf("Client-%s::Leaving Peer()", c.ID)

	var newChain bool
	for _, peer := range c.Peers {
		client, err := brpc.ConnectTo(*peer)
		if err != nil {
			log.Printf("Error while connecting to peer: %v", err)
			return err
		}
		args := &brpc.ConnectionRPC{c.ID}
		var reply brpc.BlocksRPC

		err = client.Call("NodeRPC.Peer", args, &reply)
		if err != nil {
			log.Printf("Error while peering: %v", err)
			return err
		}

		// Take longest valid chain
		if len(reply.Blocks) > len(c.blocks) {
			// Iterate over new chain.
			var invalid bool
			for i := len(reply.Blocks) - 1; i > 0; i-- {
				if reply.Blocks[i-1].Header.Hash != reply.Blocks[i].Header.PreviousBlock {
					invalid = true
					break
				}
			}
			if !invalid {
				log.Println("Incoming blockchain is valid!")
				c.blocks = reply.Blocks
				newChain = true
			}
		} else if len(c.blocks) > 0 && len(reply.Blocks) == len(c.blocks) {
			log.Println("Received chain is as large as mine")
			if reply.Blocks[len(reply.Blocks)-1].Header.Nounce > c.blocks[len(c.blocks)-1].Header.Nounce {
				log.Printf("Mine has higher nounce: %d > %d", c.blocks[len(c.blocks)-1].Header.Nounce, reply.Blocks[len(reply.Blocks)-1].Header.Nounce)
			} else {
				log.Printf("Incoming block has higher nounce, making it my main chain: %d > %d", reply.Blocks[len(reply.Blocks)-1].Header.Nounce, c.blocks[len(c.blocks)-1].Header.Nounce)
				c.blocks = reply.Blocks
				newChain = true
			}
		}
		peer.Conn = client
		log.Printf("Successfully peered with node-%s\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
	}
	if newChain {
		for _, b := range c.blocks {
			err := c.BroadcastBlock(b)
			if err != nil {
				log.Printf("Error while broadcasting newly received blockchain")
				return err
			}
		}
	}
	return nil
}

func (c *Client) GetBlocks() []node.Block {
	log.Printf("Client-%s::Entering GetBlocks()", c.ID)
	defer log.Printf("Client-%s::Leaving GetBlocks()", c.ID)

	return c.blocks
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

// Connect node to Anchor Miner
func (c *Client) Connect(host string, anchorPort string) error {
	log.Printf("Client-%s::Entering Connect()", c.ID)
	defer log.Printf("Client-%s::Leaving Connect()", c.ID)

	anchorPeer := &node.Peer{
		Host: fmt.Sprintf(host),
		Port: anchorPort}

	// Connect to Anchor Miner
	client, err := brpc.ConnectTo(*anchorPeer)
	if err != nil {
		log.Printf("Error while connecting to anchor")
		return err
	}

	// Keep connection
	c.Peers = make([]*node.Peer, 0)
	anchorPeer.Conn = client
	c.Peers = append(c.Peers, anchorPeer)

	// Request Miner to create incoming connection
	var reply int

	me := brpc.PeerRPC{
		ConnectionRPC: brpc.ConnectionRPC{PeerID: c.ID},
		Host:          c.Host,
		Port:          c.ID}

	err = client.Call("NodeRPC.Connect", &me, &reply)
	if err != nil {
		log.Printf("Error while requesting connection to anchor : %s", err)
		c.Peers = nil
		return err
	}

	var replyBlock brpc.BlocksRPC
	var sendBlock brpc.GetBlocksRPC

	err = client.Call("NodeRPC.GetBlocks", sendBlock, &replyBlock)
	if err != nil {
		log.Printf("Error while fetching block from anchor : %s", err)
		var reply int
		client.Call("NodeRPC.Disconnect", &c.ID, &reply)
		c.Peers = nil
		return err
	}
	c.blocks = replyBlock.Blocks
	for _, block := range c.blocks {
		c.ParseBlock(block)
	}
	// Restart Message Loop
	if c.trafficGeneration == true {
		go c.StartMessageLoop()
	}

	return nil
}

// Ask all node to close connection and erase current connection slice
func (c *Client) Disconnect() error {
	log.Printf("Client-%s::Entering Disconnect()", c.ID)
	defer log.Printf("Client-%s::Leaving Disconnect()", c.ID)

	for _, conn := range c.Peers {
		var reply int

		if conn.Conn == nil {
			log.Println("Error: Peer's connection is nil")
			return fmt.Errorf("Error: Peer's connection is nil")
		}
		err := conn.Conn.Call("NodeRPC.Disconnect", &c.ID, &reply)
		if err != nil {
			log.Printf("Error while disconnecting: %v", err)
			return nil
		}
	}
	c.Peers = nil

	if c.trafficGeneration == true && c.msgLoopRunning == true {
		c.msgLoopChan <- " "
		c.msgLoopRunning = false
	}

	return nil
}

// Ignore all request for bidirectionnal connection
func (c *Client) OpenConnection(host string, connectingPort string) error {
	log.Printf("Client-%s::Entering OpenConnection()", c.ID)
	defer log.Printf("Client-%s::Leaving OpenConnection()", c.ID)

	// Do we already have a connection to this peer?
	for _, p := range c.Peers {
		if p.Host == host && p.Port == connectingPort {
			log.Printf("We already have peer \"%s:%s\" in our list", host, connectingPort)
			return nil
		}
	}

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

	tempPeers := append(c.Peers, anchorPeer)
	c.Peers = tempPeers

	return nil
}

// Close connection of requesting peer
func (c *Client) CloseConnection(disconnectingPeer string) error {
	log.Printf("Client-%s::Entering CloseConnection()", c.ID)
	defer log.Printf("Client-%s::Leaving CloseConnection()", c.ID)

	for i := 0; i < len(c.Peers); i++ {
		if c.Peers[i].Port == disconnectingPeer {
			c.Peers[i] = c.Peers[len(c.Peers)-1]
			c.Peers = c.Peers[:len(c.Peers)-1]
			break
		}
	}
	return nil
}

func (c *Client) ReceiveMessage(content, temps, peer string, messageType int) error {
	log.Printf("Client-%s::Entering ReceiveMessage()", c.ID)
	defer log.Printf("Client-%s::Leaving ReceiveMessage()", c.ID)

	return nil
}

func (c *Client) HandleUiMessage(msg node.Message) error {
	log.Printf("Client-%s::Entering HandleUiMessage()", c.ID)
	defer log.Printf("Client-%s::Leaving HandleUiMessage()", c.ID)

	for _, conn := range c.Peers {
		var reply int
		message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, msg.Content, msg.Time, brpc.MessageType}

		if conn.Conn == nil {
			log.Println("Error: Peer's connection is nil")
			return fmt.Errorf("Error: Peer's connection is nil")
		}
		err := conn.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
		if err != nil {
			log.Printf("Error while delivering message: %v", err)
			log.Printf("Closing connection")
			c.CloseConnection(conn.Port)
			return nil
		}
	}
	return nil
}

func (c *Client) StartMessageLoop() error {
	log.Printf("Client-%s::Entering StartMessageLoop()", c.ID)
	defer log.Printf("Client-%s::Leaving StartMessageLoop()", c.ID)

	time.Sleep(20 * time.Second)

	c.msgLoopRunning = true
	for {
		select {
		case <-c.msgLoopChan:
			return nil
		default:
			time.Sleep(5 * time.Second)
			for _, peer := range c.Peers {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour!", time.Now().Format(time.RFC3339Nano), brpc.MessageType}

				if peer.Conn == nil {
					log.Println("Error: Peer's connection is nil")
					return fmt.Errorf("Error: Peer's connection is nil")
				}
				err := peer.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
				if err != nil {
					log.Printf("Error while trying to deliver message: %v", err)
					log.Printf("Closing connection")
					c.CloseConnection(peer.Port)
					return nil
				}

			}
		}
	}

	return nil
}

func (c *Client) ParseBlock(block node.Block) error {
	if c.uiChannel != nil {
		for _, msg := range block.Messages {
			log.Printf("Filling channel...")
			c.uiChannel <- node.Message{msg.Peer, msg.Content, msg.Time}
		}

	}
	return nil
}

func (c *Client) BroadcastBlock(b node.Block) error {
	log.Printf("Client-%s::Entering BroadcastBlock()", c.ID)
	defer log.Printf("Client-%s::Leaving BroadcastBlock()", c.ID)

	if len(c.Peers) == 0 {
		log.Println("No peers are defined. Exiting.")
		return nil
	}

	for _, peer := range c.Peers {
		if peer.Conn == nil {
			log.Printf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))

		}
		args := brpc.BlockRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: c.ID},
			Block:         b}
		var reply *int

		if peer.Conn == nil {
			log.Println("Error: Peer's connection is nil")
			return fmt.Errorf("Error: Peer's connection is nil")
		}
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		if err != nil {
			log.Printf("Error while delivering block: %v", err)
			return err
		}
	}
	return nil
}

func (c *Client) ReceiveBlock(block node.Block, peer string) error {
	log.Printf("Client-%s::Entering ReceiveBlock()", c.ID)
	defer log.Printf("Client-%s::Leaving ReceiveBlock()", c.ID)

	log.Printf("Received block %v from \"%s\"", block, peer)
	// Do we already have this block in the chain?
	for _, b := range c.blocks {
		if reflect.DeepEqual(b, block) {
			log.Println("We already have this block in the chain. Discarding block")
			return nil
		}
	}

	if len(c.blocks) > 0 && (block.Header.PreviousBlock != c.blocks[len(c.blocks)-1].Header.Hash) {
		log.Println("Incoming block does not point at the head of the chain")
		return nil
	}

	valid := true

	// We don't have any blocks and the received one is a genesis block
	// if len(c.blocks) == 0 && block.Header.PreviousBlock == ([sha256.Size]byte{}) && ) {
	// 	valid = false
	// }

	//que le hash est correct (bonne difficulter)
	if valid {
		header := node.Header{
			PreviousBlock: block.Header.PreviousBlock,
			Nounce:        block.Header.Nounce,
			Date:          block.Header.Date,
			ContentHash:   block.Header.ContentHash,
		}
		header.Hash = [sha256.Size]byte{}
		hash := sha256.Sum256([]byte(fmt.Sprintf("%v", header)))

		if block.Header.Hash != hash {
			log.Printf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
			return fmt.Errorf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
		}
		contentHash := sha256.Sum256([]byte(fmt.Sprintf("%v", block.Messages)))
		if block.Header.ContentHash != contentHash {
			log.Printf("Error while validating content hash! (%v != %v)", contentHash, block.Header.ContentHash)
			return fmt.Errorf("Error while validating content hash! (%v != %v)", contentHash, block.Header.ContentHash)
		}

		firstCharacters := string(hash[:node.MiningDifficulty])
		if strings.Count(firstCharacters, "0") == node.MiningDifficulty && hash == block.Header.Hash {
			log.Printf("Locking blocksMutex mutex")
			c.blocksMutex.Lock()
			log.Printf("Locked blocksMutex mutex")
			log.Println("Received block's hash is valid!")
			log.Println("But we still need to do some checking")

		} else {
			log.Printf("Received block's hash is not valid! Discarding block (%v != %v)", block.Header.Hash, hash)
			return nil
		}

	} else {
		log.Println("Block does not exist in my chain! Ask my peers for the missing range of blocks")
		// TODO: Implement fetching range of blocks from peers
		return nil
	}

	//TODO: add timestamp gap checking between messages

	log.Println("Appending block to the chain")
	// TODO: Use mutex?
	tempBlocks := append(c.blocks, block)
	c.blocks = tempBlocks
	c.blocksMutex.Unlock()
	log.Printf("Unlocked blocksMutex mutex")

	c.ParseBlock(block)

	return c.BroadcastBlock(block)
}
