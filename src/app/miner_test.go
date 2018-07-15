package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"crypto/sha256"
	"reflect"
	"sync"
	"testing"
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
				make(chan node.Message, 100),
				make(chan node.Block, 10),
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
			m.mining()

		})
	}
}

func TestMiner_findingNounce(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []string
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan Message
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
				make(chan Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
				sync.Mutex{},
			},
			args{
				
			}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miner.findingNounce() = %v, want %v", got, tt.want)
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
			if got := m.CreateBlock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miner.CreateBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
