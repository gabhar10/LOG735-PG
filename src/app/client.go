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
	}

	PeerConnection struct {
		ID			string
		conn		*rpc.Client
	}
)

func NewClient(port, peers string) node.Node {
	c := &Client{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		strings.Split(peers, " "),
		new(brpc.NodeRPC),
		make([]PeerConnection,0),
	}
	c.rpcHandler.Node = c
	return c
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

func (c Client) ReceiveMessage(content string, hello time.Time) {

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
	}


	// start trafic generation
	go c.StartMessageLoop()

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

func (c Client) StartMessageLoop() error{
	log.Printf("Starting message loop")
	for{
		for _, conn := range c.connections {
			time.Sleep(5 * time.Second)
			log.Printf("CLIENT : Sending message to %s\n", conn.ID)
			var reply int
			message := brpc.MessageRPC{brpc.ConnectionRPC{c.ID}, "Bonjour", time.Now()}
			conn.conn.Call("NodeRPC.DeliverMessage", message, &reply)
		}
	}

	return nil
}
