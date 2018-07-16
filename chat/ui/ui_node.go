package ui


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
	Ui_Node struct {
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

func NewUiNode(port, peers string) node.Node {
	c := &Ui_Node{
		port,
		make([]node.Block, node.MinBlocksReturnSize),
		strings.Split(peers, " "),
		new(brpc.NodeRPC),
		make([]PeerConnection,0),
	}
	c.rpcHandler.Node = c
	return c
}

func (c *Ui_Node) SetupRPC(port string) error {
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

func (c Ui_Node) ReceiveMessage(content string, hello time.Time) {

}

func (c Ui_Node) Peer() error {
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

	return nil
}

func (c Ui_Node) GetBlocks() []node.Block {
	return nil
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

func (c Ui_Node) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// CLIENT-04, CLIENT-06
	// To implement
	return nil
}

func (c Ui_Node) Disconnect() error {
	// Disconnect (RPC) to peers
	// CLIENT-05
	// To implement
	return nil
}
