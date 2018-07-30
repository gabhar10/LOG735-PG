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
	"time"
)

type (
	Client struct {
		ID             string            // i.e. Run-time port associated to container
		blocks         []node.Block      // Can be a subset of the full chain
		Peers          []*node.Peer      // Slice of Peers
		rpcHandler     *brpc.NodeRPC     // Handler for RPC calls
		uiChannel      chan node.Message // Channel to send message from node to chat application
		nodeChannel    chan node.Message // Channel to send message from chat application to node
		msgLoopChan    chan string
		msgLoopRunning bool
	}
)

func NewClient(port string, peers []*node.Peer, uiChannel chan node.Message, nodeChannel chan node.Message) node.Node {
	c := &Client{
		port,
		[]node.Block{},
		peers,
		new(brpc.NodeRPC),
		uiChannel,
		nodeChannel,
		make(chan string),
		false,
	}
	c.rpcHandler.Node = c
	return c
}

func (c *Client) Start() {
	go c.StartMessageLoop()
}

func (c *Client) SetupRPC() error {
	log.Println("Entering SetupRPC()")
	defer log.Println("Leaving SetupRPC()")

	s := rpc.NewServer()
	s.Register(c.rpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", c.ID))
	if err != nil {
		return err
	}
	go s.Accept(listener)
	return nil
}

func (c *Client) Peer() error {
	for _, peer := range c.Peers {
		client, err := brpc.ConnectTo(*peer)
		if err != nil {
			return err
		}
		args := &brpc.ConnectionRPC{c.ID}
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

func (c *Client) GetBlocks() []node.Block {
	return c.blocks
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

// Connect node to Anchor Miner
func (c *Client) Connect(anchorPort string) error {

	anchorPeer := &node.Peer{
		Host: fmt.Sprintf("node-%s", anchorPort),
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
	err = client.Call("NodeRPC.Connect", c.ID, &reply)
	if err != nil {
		log.Printf("Error while requesting connection to anchor : %s", err)
		c.Peers = nil
		return err
	}

	// Restart Message Loop
	go c.StartMessageLoop()

	return nil
}

// Ask all node to close connection and erase current connection slice
func (c *Client) Disconnect() error {

	log.Printf("Disconnecting from network")
	for _, conn := range c.Peers {
		var reply int
		err := conn.Conn.Call("NodeRPC.Disconnect", &c.ID, &reply)
		if err != nil {
			return err
		}
	}
	c.Peers = nil
	if c.msgLoopRunning == true {
		c.msgLoopChan <- " "
		c.msgLoopRunning = false
	}

	return nil
}

// Ignore all request for bidirectionnal connection
func (c *Client) OpenConnection(connectingPort string) error {
	return nil
}

// Close connection of requesting peer
func (c *Client) CloseConnection(disconnectingPeer string) error {
	for i := 0; i < len(c.Peers); i++ {
		if c.Peers[i].Port == disconnectingPeer {
			c.Peers[i].Conn.Close()
			c.Peers[i] = c.Peers[len(c.Peers)-1]
			c.Peers = c.Peers[:len(c.Peers)-1]
			break
		}
	}
	return nil
}

func (c *Client) ReceiveMessage(content string, temps time.Time, peer string, messageType int) {
	if c.uiChannel != nil {
		c.uiChannel <- node.Message{peer, content, temps}
	}
}

func (c *Client) HandleUiMessage(msg node.Message) error {
	for _, conn := range c.Peers {
		var reply int
		message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, msg.Content, msg.Time, brpc.MessageType}
		err := conn.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) StartMessageLoop() error {
	c.msgLoopRunning = true
	for {
		select {
		case <-c.msgLoopChan:
			return nil
		default:
			time.Sleep(5 * time.Second)
			for _, peer := range c.Peers {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour", time.Now(), brpc.MessageType}
				peer.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
			}
			if c.nodeChannel != nil {
				c.uiChannel <- node.Message{c.ID, "Bonjour", time.Now()}
			}
		}

	}

	return nil
}

func (c *Client) BroadcastBlock(b node.Block) error {
	log.Println("Entering BroadcastBlock()")
	defer log.Println("Leaving BroadcastBlock()")

	if len(c.Peers) == 0 {
		log.Println("No peers are defined. Exiting.")
		return nil
	}

	for _, peer := range c.Peers {
		if peer.Conn == nil {
			return fmt.Errorf("RPC connection handler of peer %s is nil", fmt.Sprintf("%s:%s", peer.Host, peer.Port))

		}
		args := brpc.BlockRPC{
			ConnectionRPC: brpc.ConnectionRPC{PeerID: c.ID},
			Block:         b}
		var reply *int
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ReceiveBlock(block node.Block) error {
	log.Println("Client::Entering ReceiveBlock()")
	defer log.Println("Client::Leaving ReceiveBlock()")

	log.Printf("Block header hash: %v, Block header nounce: %v, Block header date: %v", block.Header.Hash, block.Header.Nounce, block.Header.Date)
	// Do we already have this block in the chain?
	for _, b := range c.blocks {
		if reflect.DeepEqual(b, block) {
			log.Println("We already have this block in the chain. Discarding block")
			return nil
		}
	}

	valid := false
	//es-ce que le bloc precedant existe dans la chaine
	for i := len(c.blocks) - 1; i >= 0; i-- {
		if block.Header.PreviousBlock == c.blocks[i].Header.Hash {
			valid = true
			break
		}
	}

	// We don't have any blocks and the received one is a genesis block
	if block.Header.PreviousBlock == ([sha256.Size]byte{}) && len(c.blocks) == 0 {
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
			log.Println("Received block's hash is valid!")
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

	return c.BroadcastBlock(block)
}
