package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"time"
)

type (
	Client struct {
		ID         string       // i.e. Run-time port associated to container
		blocks     []node.Block // Can be a subset of the full chain
		peers      []string     // Slice of IDs
		rpcHandler *brpc.NodeRPC
		connections []PeerConnection
		uiChannel chan node.Message
		nodeChannel chan node.Message
	}

	PeerConnection struct {
		ID   string
		conn *rpc.Client
	}
)

func NewClient(port, peers string, uiChannel chan node.Message, nodeChannel chan node.Message) node.Node {
	c := &Client{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		strings.Split(peers, " "),
		new(brpc.NodeRPC),
		make([]PeerConnection,0),
		uiChannel,
		nodeChannel,
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
	log.Printf("Listening on TCP port %s\n", port)
	go http.Serve(l, nil)
	return nil
}


func (c Client) Peer() error {
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
		if reply.Blocks != nil {
			return fmt.Errorf("Blocks are not defined")
		}
		var newConnection = new(PeerConnection)
		newConnection.ID = peer
		newConnection.conn = client
		c.connections = append(c.connections, *newConnection)
		log.Printf("Successfully peered with node-%s\n", peer)
	}


	// start trafic generation
	go c.StartMessageLoop()
	go c.HandleUiMessage()
	return nil
}

func (c Client) GetBlocks() []node.Block {
	return c.blocks
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

func (c Client) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// CLIENT-04, CLIENT-06
	// To implement
	return nil
}

func (c Client) Disconnect() error {
	// Disconnect (RPC) to peers
	// CLIENT-05
	// To implement
	return nil
}

func (c Client) ReceiveMessage(content string, temps time.Time, peer string) {
	if c.uiChannel != nil {
		c.uiChannel <- node.Message{peer, content, temps}
	}
}

func (c Client) HandleUiMessage() error{
	if c.nodeChannel != nil {
		for {
			var msg node.Message
			msg = <-c.nodeChannel
			for _, conn := range c.connections {
				var reply int
				message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, msg.Content, time.Now()}
				conn.conn.Call("NodeRPC.DeliverMessage", message, &reply)
			}
		}
	}
	return nil
}

func (c Client) StartMessageLoop() error{
	log.Printf("Starting message loop")
	for{
		time.Sleep(5 * time.Second)
		for _, conn := range c.connections {
			var reply int
			message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour", time.Now()}
			conn.conn.Call("NodeRPC.DeliverMessage", message, &reply)
		}
		if c.nodeChannel != nil {
			c.uiChannel <- node.Message{c.ID, "Bonjour", time.Now()}
		}
	}

	return nil
}


func (c Client) ReceiveBlock(block node.Block) {
	// Do nothing
}
