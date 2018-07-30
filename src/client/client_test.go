package client

import (
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"testing"
	"time"
)

func TestClient_Peer(t *testing.T) {
	type fields struct {
		Host		string
		ID          string
		blocks      []node.Block
		peers       []*node.Peer
		rpcHandler  *brpc.NodeRPC
		uiChannel   chan node.Message
		nodeChannel chan node.Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Peer with miner",
			fields: fields{
				Host : "127.0.0.1",
				ID: "9001",
				peers: func() []*node.Peer {
					driver := node.Peer{Host: "127.0.0.1", Port: "9001"}
					m := miner.NewMiner("9002", []*node.Peer{&driver}).(*miner.Miner)
					err := m.SetupRPC()
					if err != nil {
						t.Errorf("Error while trying to setup RPC: %v", err)
					}
					return []*node.Peer{&node.Peer{Host: "127.0.0.1", Port: "9002"}}
				}(),
				uiChannel:   nil,
				nodeChannel: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.Host, tt.fields.ID, tt.fields.peers, tt.fields.uiChannel, tt.fields.nodeChannel)
			if err := c.Peer(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Peer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ReceiveBlock(t *testing.T) {
	type fields struct {
		Host		   string
		ID             string
		blocks         []node.Block
		Peers          []*node.Peer
		rpcHandler     *brpc.NodeRPC
		uiChannel      chan node.Message
		nodeChannel    chan node.Message
		msgLoopChan    chan string
		msgLoopRunning bool
	}
	type args struct {
		block node.Block
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Sunny day",
			fields: fields{
				Host:		 "127.0.0.1",
				ID:          "9002",
				Peers:       []*node.Peer{},
				uiChannel:   nil,
				nodeChannel: make(chan node.Message, node.MessagesChannelSize),
			},
			args: args{
				block: func() node.Block {
					m := miner.NewMiner("8888", []*node.Peer{}).(*miner.Miner)
					messages := [node.BlockSize]node.Message{}
					for i := 0; i < node.BlockSize; i++ {
						messages[i] = node.Message{
							Peer:    "9000",
							Content: "Salut!",
							Time:    time.Now().Format(time.RFC3339),
						}
					}

					for _, msg := range messages {
						m.IncomingMsgChan <- msg
					}
					m.Start()
					// Wait for miner to mine a block
					found := false
					for i := 0; i < 5; i++ {
						if len(m.GetBlocks()) > 0 {
							found = true
							break
						} else {
							time.Sleep(time.Second * 2)
						}
					}
					if !found {
						t.Errorf("Miner could not mine a block")
					}
					return m.GetBlocks()[0]
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.Host ,tt.fields.ID, tt.fields.Peers, tt.fields.uiChannel, tt.fields.nodeChannel)
			if err := c.ReceiveBlock(tt.args.block); (err != nil) != tt.wantErr {
				t.Errorf("Client.ReceiveBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
