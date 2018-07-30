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
		Host		   string
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

func NewClient(host string, port string, peers []*node.Peer, uiChannel chan node.Message, nodeChannel chan node.Message) node.Node {
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
	}
	c.rpcHandler.Node = c
	return c
}

func (c *Client) Start() {
	log.Printf("Client-%s::Entering Start()", c.ID)
	defer log.Printf("Client-%s::Leaving Start()", c.ID)

	go c.StartMessageLoop()
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
		peer.Conn = client
		log.Printf("Successfully peered with node-%s\n", fmt.Sprintf("%s:%s", peer.Host, peer.Port))
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
		Host: c.Host,
		Port: c.ID}


	err = client.Call("NodeRPC.Connect", &me, &reply)
	if err != nil {
		log.Printf("Error while requesting connection to anchor : %s", err)
		c.Peers = nil
		return err
	}

	// Restart Message Loop
	//go c.StartMessageLoop()

	return nil
}

// Ask all node to close connection and erase current connection slice
func (c *Client) Disconnect() error {
	log.Printf("Client-%s::Entering Disconnect()", c.ID)
	defer log.Printf("Client-%s::Leaving Disconnect()", c.ID)

	for _, conn := range c.Peers {
		var reply int
		err := conn.Conn.Call("NodeRPC.Disconnect", &c.ID, &reply)
		if err != nil {
			log.Printf("Error while disconnecting: %v", err)
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
func (c *Client) OpenConnection(host string, connectingPort string) error {
	log.Printf("Client-%s::Entering OpenConnection()", c.ID)
	defer log.Printf("Client-%s::Leaving OpenConnection()", c.ID)

	return nil
}

// Close connection of requesting peer
func (c *Client) CloseConnection(disconnectingPeer string) error {
	log.Printf("Client-%s::Entering CloseConnection()", c.ID)
	defer log.Printf("Client-%s::Leaving CloseConnection()", c.ID)

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
		err := conn.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
		if err != nil {
			log.Printf("Error while delivering message: %v", err)
			return err
		}
	}
	return nil
}

func (c *Client) StartMessageLoop() error {
	log.Printf("Client-%s::Entering StartMessageLoop()", c.ID)
	defer log.Printf("Client-%s::Leaving StartMessageLoop()", c.ID)

	c.msgLoopRunning = true
	for {
		select {
		case <-c.msgLoopChan:
			return nil
		default:
			time.Sleep(5 * time.Second)
			for _, peer := range c.Peers {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour", time.Now().Format(time.RFC3339Nano), brpc.MessageType}
				err := peer.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
				if err != nil {
					log.Printf("Error while trying to deliver message: %v", err)
					return fmt.Errorf("Error while trying to deliver message: %v", err)
				}

			}
			if c.nodeChannel != nil {
				c.uiChannel <- node.Message{c.ID, "Bonjour", time.Now().Format(time.RFC3339Nano)}
			}
		}
	}

	return nil
}

func (c *Client) GetAnchorBlock() ([]node.Block, error){
	if len(c.Peers) > 0 {
		if c.Peers[0].Conn != nil{
			var reply brpc.BlocksRPC
			c.Peers[0].Conn.Call("NodeRPC.GetBlocks", nil ,&reply)
			return reply.Blocks, nil
		} else {
			return nil, fmt.Errorf("no open connection with peer")
		}
	} else {
		return nil, fmt.Errorf("there are no peer available to give a block")
	}
}


func (c *Client) ParseBlock(block node.Block) (error){
	if c.uiChannel != nil {
		for _, msg := range block.Messages{
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
		err := peer.Conn.Call("NodeRPC.DeliverBlock", &args, &reply)
		if err != nil {
			log.Printf("Error while delivering block: %v", err)
			return err
		}
	}
	return nil
}

func (c *Client) ReceiveBlock(block node.Block) error {
	log.Printf("Client-%s::Entering ReceiveBlock()", c.ID)
	defer log.Printf("Client-%s::Leaving ReceiveBlock()", c.ID)

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

		if block.Header.Hash != hash {
			log.Printf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
			return fmt.Errorf("Error while validating hash! (%v != %v)", hash, block.Header.Hash)
		}

		firstCharacters := string(hash[:node.MiningDifficulty])
		if strings.Count(firstCharacters, "0") == node.MiningDifficulty && hash == block.Header.Hash {
			log.Println("Check if blocks are identical")
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
	c.ParseBlock(block)

	return c.BroadcastBlock(block)
}
