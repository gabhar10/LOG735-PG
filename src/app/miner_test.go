package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"fmt"
	"net"
	"net/rpc"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMiner_findingNounce(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	type args struct {
		block *node.Block
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [sha256.Size]byte
		wantErr bool
	}{
		{
			name: "sunny day",
			fields: fields{
				"1",
				make([]node.Block, 2),
				make([]node.Peer, 1),
				nil,
				make(chan node.Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
				new(sync.Mutex),
			},
			args: args{
				func() *node.Block {
					messages := [node.BlockSize]node.Message{}
					for i := 0; i < node.BlockSize; i++ {
						messages[i] = node.Message{"Salut!", time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
					}
					return &node.Block{Header: node.Header{PreviousBlock: [sha256.Size]byte{}, Date: time.Now()}, Messages: messages}
				}(),
			},
			want:    [sha256.Size]byte{},
			wantErr: false,
		},
		{
			name: "Use quit channel",
			fields: fields{
				quit: func() chan bool {
					c := make(chan bool, 1)
					c <- true
					return c
				}(),
			},
			want:    [sha256.Size]byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			got, err := m.findingNounce(tt.args.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("Miner.findingNounce() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				firstCharacters := string(got[:node.MiningDifficulty])
				if strings.Count(firstCharacters, "0") != node.MiningDifficulty {
					t.Errorf("first %v characters are not zeros", node.MiningDifficulty)
					return
				}
			}
		})
	}
}

func TestMiner_CreateBlock(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   node.Block
	}{
		{
			"sunny day",
			fields{
				"1",
				make([]node.Block, 2),
				make([]node.Peer, 1),
				nil,
				make(chan node.Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
				new(sync.Mutex),
			},
			func() node.Block {
				messages := [node.BlockSize]node.Message{}
				for i := 0; i < node.BlockSize; i++ {
					messages[i] = node.Message{"Salut!", time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
				}
				return node.Block{Messages: messages}
			}(),
		},
	}
	for _, tt := range tests {
		for i := 0; i < node.BlockSize; i++ {
			tt.fields.incomingMsgChan <- node.Message{"Salut!", time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
		}
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			got := m.CreateBlock()
			if !reflect.DeepEqual(got.Messages, tt.want.Messages) {
				t.Errorf("Miner.CreateBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMiner(t *testing.T) {
	type args struct {
		port  string
		peers []node.Peer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sunny day",
			args: args{
				port: "123",
				peers: []node.Peer{
					node.Peer{
						Host: "456",
						Port: "123",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMiner(tt.args.port, tt.args.peers)
			miner, ok := got.(*Miner)
			if ok {
				if miner.ID != tt.args.port || len(miner.peers) != 1 || miner.peers[0] != tt.args.peers[0] {
					t.Errorf("NewMiner() did not populate with port %s and peers %s", tt.args.port, tt.args.peers)
				}
				if cap(miner.blocks) != node.MinBlocksReturnSize || cap(miner.incomingMsgChan) != node.MessagesChannelSize || cap(miner.incomingBlockChan) != node.BlocksChannelSize {
					t.Errorf("NewMiner() did not return a miner with attributes of proper length")
				}
			} else {
				t.Errorf("NewMiner() did not return a structure of type *Miner")
			}
		})
	}
}

func TestMiner_Start(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Sunny Day",
			fields: fields{
				mutex: new(sync.Mutex),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			m.Start()
		})
	}
}

func TestMiner_SetupRPC(t *testing.T) {
	var l net.Listener
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	type args struct {
		port string
	}
	tests := []struct {
		name     string
		fields   fields
		preFunc  func()
		postFunc func()
		args     args
		wantErr  bool
	}{
		{
			name:     "Port is already taken",
			preFunc:  func() { l, _ = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", "8888")) },
			postFunc: func() { l.Close() },
			fields: fields{
				mutex: new(sync.Mutex),
			},
			args: args{
				port: "8888",
			},
			wantErr: true,
		},
		{
			name:     "Sunny day",
			preFunc:  func() {},
			postFunc: func() {},
			fields: fields{
				mutex: new(sync.Mutex),
			},
			args: args{
				port: "8888",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			m := NewMiner(tt.args.port, []node.Peer{})
			if err := m.SetupRPC(tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Miner.SetupRPC() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.postFunc()
		})
	}
}

func TestMiner_Peer(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Sunny day",
			fields: fields{
				ID: "9001",
				peers: func() []node.Peer {
					driver := node.Peer{Host: "127.0.0.1", Port: "9001"}
					m := NewMiner("9002", []node.Peer{driver}).(*Miner)
					// Given an RPC handler already exist from previous test, use it
					if rpc.Register(m.rpcHandler) != nil {
						return []node.Peer{node.Peer{Host: "127.0.0.1", Port: "8888"}}
					}
					m.SetupRPC("9002")
					return []node.Peer{node.Peer{Host: "127.0.0.1", Port: "9002"}}
				}(),
				mutex: new(sync.Mutex),
			},
			wantErr: false,
		},
		{
			name: "Can't connect",
			fields: fields{
				ID:    "9001",
				peers: []node.Peer{node.Peer{Host: "127.0.0.1", Port: "9003"}},
				mutex: new(sync.Mutex),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiner(tt.fields.ID, tt.fields.peers)
			if err := m.Peer(); (err != nil) != tt.wantErr {
				t.Errorf("Miner.Peer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiner_Broadcast(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Sunny day",
			fields: fields{
				ID:    "8999",
				peers: []node.Peer{},
				mutex: new(sync.Mutex),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiner(tt.fields.ID, tt.fields.peers).(*Miner)
			if err := m.Broadcast(); (err != nil) != tt.wantErr {
				t.Errorf("Miner.Broadcast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiner_ReceiveMessage(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	type args struct {
		content string
		temps   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Sunny day",
			fields: fields{
				incomingMsgChan: make(chan node.Message, 1),
				mutex:           new(sync.Mutex),
			},
			args: args{
				content: "Hello",
				temps:   time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			m.ReceiveMessage(tt.args.content, tt.args.temps)
			mes := <-m.incomingMsgChan
			if mes.Content != tt.args.content || mes.Time.After(time.Now()) {
				t.Errorf("Miner.ReceiveMessage() did not return expected values: ")
			}
		})
	}
}

func TestMiner_ReceiveBlock(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	type args struct {
		block node.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Sunny day",
			fields: fields{
				quit:  make(chan bool, 1),
				mutex: new(sync.Mutex),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			m.ReceiveBlock(tt.args.block)
		})
	}
}

func TestMiner_mining(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []node.Peer
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   node.Block
	}{
		{
			name: "sunny day",
			fields: fields{
				"1",
				make([]node.Block, 2),
				make([]node.Peer, 1),
				nil,
				make(chan node.Message, node.MessagesChannelSize),
				make(chan node.Block, node.BlocksChannelSize),
				make(chan bool),
				new(sync.Mutex),
			},
			want: func() node.Block {
				messages := [node.BlockSize]node.Message{}
				for i := 0; i < node.BlockSize; i++ {
					messages[i] = node.Message{Content: "Salut!"}
				}
				return node.Block{Messages: messages}
			}(),
		},
		{
			name: "Quit mining",
			fields: fields{
				"1",
				make([]node.Block, 2),
				make([]node.Peer, 1),
				nil,
				make(chan node.Message, node.MessagesChannelSize),
				make(chan node.Block, node.BlocksChannelSize),
				func() chan bool {
					c := make(chan bool, 1)
					c <- true
					return c
				}(),
				new(sync.Mutex),
			},
			want: node.Block{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:                tt.fields.ID,
				blocks:            tt.fields.blocks,
				peers:             tt.fields.peers,
				rpcHandler:        tt.fields.rpcHandler,
				incomingMsgChan:   tt.fields.incomingMsgChan,
				incomingBlockChan: tt.fields.incomingBlockChan,
				quit:              tt.fields.quit,
				mutex:             tt.fields.mutex,
			}
			// Send messages in go routine
			go func() {
				for i := 0; i < node.BlockSize; i++ {
					m.incomingMsgChan <- node.Message{
						Content: "Salut!",
						Time:    time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
				}
			}()

			got := m.mining()
			for _, v := range tt.want.Messages {
				for _, w := range got.Messages {
					if v.Content != w.Content {
						t.Errorf("Received messages (%v) != expected messages (%v)", v, w)
					}
				}
			}
		})
	}
}
