package client

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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
		make([]node.Block, node.MinBlocksReturnSize),
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

func (c *Client) SetupRPC(port string) error {
	rpc.Register(c.rpcHandler)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	go http.Serve(l, nil)
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
		if len(reply.Blocks) < node.MinBlocksReturnSize {
			return fmt.Errorf("Returned size of blocks is below %d", node.MinBlocksReturnSize)
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
	//go c.StartMessageLoop()

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

}

func (c *Client) HandleUiMessage(message node.Message) error {

	log.Printf("Entering HandleUiMessage()")
	defer log.Printf("Leaving HandleUiMessage()")
	for _, conn := range c.Peers {
		var reply int
		message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, message.Content, message.Time, brpc.MessageType}
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


func (c *Client) ParseChain([]node.Block) (error){
	return nil
}

func (c *Client) ReceiveBlock(block node.Block) {
	if c.uiChannel != nil {
		for _, msg := range block.Messages{
			log.Printf("Filling channel...")
			c.uiChannel <- node.Message{msg.Peer, msg.Content, msg.Time}
		}

	}
}
