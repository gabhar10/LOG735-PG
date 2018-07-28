package client

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"
	"log"
)

type (
	Client struct {
		ID          string            // i.e. Run-time port associated to container
		blocks      []node.Block      // Can be a subset of the full chain
		peers       []node.Peer       // Slice of peers
		rpcHandler  *brpc.NodeRPC     // Handler for RPC calls
		connections []node.PeerConnection  // Slice of all current connection
		uiChannel   chan node.Message // Channel to send message from node to chat application
		nodeChannel chan node.Message // Channel to send message from chat application to node
		msgLoopChan chan string
		connected	bool
	}
)

func NewClient(port string,
	peers []node.Peer,
	uiChannel chan node.Message,
	nodeChannel chan node.Message) node.Node {
	c := &Client{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		peers,
		new(brpc.NodeRPC),
		make([]node.PeerConnection, 0),
		uiChannel,
		nodeChannel,
		make(chan string),
		true,
	}
	c.rpcHandler.Node = c
	return c
}

func (c *Client) Start() {

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
	for _, peer := range c.peers {
		client, err := brpc.ConnectTo(peer)
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

		var newConnection = new(node.PeerConnection)
		newConnection.ID = peer.Port
		newConnection.Conn = client
		c.connections = append(c.connections, *newConnection)
	}

	// start trafic generation
	go c.StartMessageLoop()
	go c.HandleUiMessage()
	//go c.Disconnect()
	return nil
}

func (c Client) GetBlocks() []node.Block {
	return c.blocks
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

func (c *Client) Connect(anchorPort string) error {

	log.Printf("Establishing bidirectional connection to %s", anchorPort)
	anchorPeer := node.Peer{
		Host: fmt.Sprintf("node-%s", anchorPort),
		Port: anchorPort}

	client, err := brpc.ConnectTo(anchorPeer)
	if err != nil {
		log.Printf("Error while connecting to anchor")
		return err
	}

	c.connections = make([]node.PeerConnection, 0)

	var newConnection = new(node.PeerConnection)
	newConnection.ID = anchorPeer.Port
	newConnection.Conn = client
	c.connections = append(c.connections, *newConnection)

	var reply int
	err = client.Call("NodeRPC.Connect", c.ID, &reply)
	if err != nil{
		log.Printf("Error while requesting connection to anchor : %s", err)

	}

	go c.StartMessageLoop()

	return nil
}

func (c *Client) Disconnect() error {

	if c.connected == true{
		// Disconnect (RPC) to peers
		// CLIENT-05
		for _, conn := range c.connections {
			var reply int
			conn.Conn.Call("NodeRPC.Disconnect", c.ID, &reply)

		}
		c.connections = nil
		c.msgLoopChan <- " "
		log.Printf("Disconnecting from network")
	}

	return nil
}

func (c *Client) CloseConnection(disconnectingPeer string) error{
	for i := 0; i < len(c.connections); i++{
		if c.connections[i].ID == disconnectingPeer{
			c.connections[i].Conn.Close()
			c.connections[i] = c.connections[len(c.connections)-1]
			c.connections = c.connections[:len(c.connections)-1]
			break
		}
	}
	return nil
}

func (c *Client) OpenConnection(connectingPort string) error{
	return nil
}

func (c Client) ReceiveMessage(content string, temps time.Time, peer string) {
	if c.uiChannel != nil {
		c.uiChannel <- node.Message{peer, content, temps}
	}
}

func (c Client) HandleUiMessage() error {
	if c.nodeChannel != nil {
		for {
			var msg node.Message
			msg = <-c.nodeChannel
			for _, conn := range c.connections {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, msg.Content, time.Now()}
				conn.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
			}
		}
	}
	return nil
}

func (c Client) StartMessageLoop() error {
	for {
		select{
		case <- c.msgLoopChan:
			return nil
		default:
			time.Sleep(5 * time.Second)
			for _, conn := range c.connections {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour", time.Now()}
				conn.Conn.Call("NodeRPC.DeliverMessage", message, &reply)
			}
			if c.nodeChannel != nil {
				c.uiChannel <- node.Message{c.ID, "Bonjour", time.Now()}
			}
		}
	}

	return nil
}

func (c Client) ReceiveBlock(block node.Block) {
	// Do nothing
}
