package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMiner_mining(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []string
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"sunny day",
			fields{
				"1",
				make([]node.Block, 2),
				make([]string, 1),
				nil,
				make(chan node.Message, node.MessagesChannelSize),
				make(chan node.Block, node.BlocksChannelSize),
				make(chan bool),
				sync.Mutex{},
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

			// Send messages in go routine
			go func() {
				for i := 0; i < node.BlockSize; i++ {
					m.incomingMsgChan <- node.Message{
						Content: "Salut!",
						Time:    time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
				}
			}()

			beforeSize := len(m.blocks)
			m.mining()
			if beforeSize == len(m.blocks) {
				t.Errorf("mining() did not create a new block")
			}
		})
	}
}

func TestMiner_findingNounce(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []string
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             sync.Mutex
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
			"sunny day",
			fields{
				"1",
				make([]node.Block, 2),
				make([]string, 1),
				nil,
				make(chan node.Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
				sync.Mutex{},
			},
			args{
				func() *node.Block {
					messages := [node.BlockSize]node.Message{}
					for i := 0; i < node.BlockSize; i++ {
						messages[i] = node.Message{"Salut!", time.Date(2018, 7, 15, 8, 0, 0, 0, time.UTC)}
					}
					return &node.Block{Header: node.Header{PreviousBlock: [sha256.Size]byte{}, Date: time.Now()}, Messages: messages}
				}(),
			},
			[sha256.Size]byte{},
			false,
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
				return
			}
			firstCharacters := string(got[:node.MiningDifficulty])

			if strings.Count(firstCharacters, "0") != node.MiningDifficulty {
				t.Errorf("first %v characters are not zeros", node.MiningDifficulty)
				return
			}
		})
	}
}

func TestMiner_CreateBlock(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []string
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan node.Message
		incomingBlockChan chan node.Block
		quit              chan bool
		mutex             sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   node.Block
	}{
		// TODO: Add test cases.
		{
			"sunny day",
			fields{
				"1",
				make([]node.Block, 2),
				make([]string, 1),
				nil,
				make(chan node.Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
				sync.Mutex{},
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
