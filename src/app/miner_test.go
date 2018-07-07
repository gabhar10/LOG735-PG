package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"testing"
)

func TestMiner_mining(t *testing.T) {
	type fields struct {
		ID              string
		blocks          []node.Block
		peers           []string
		rpcHandler      *brpc.NodeRPC
		messageQueue    []Message
		messageChannels []chan Message
		blocksChannels  []chan node.Block
		miningBlock     node.Block
		quit            chan bool
	}
	type args struct {
		block      *node.Block
		difficulty int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				ID:              tt.fields.ID,
				blocks:          tt.fields.blocks,
				peers:           tt.fields.peers,
				rpcHandler:      tt.fields.rpcHandler,
				messageQueue:    tt.fields.messageQueue,
				messageChannels: tt.fields.messageChannels,
				blocksChannels:  tt.fields.blocksChannels,
				miningBlock:     tt.fields.miningBlock,
				quit:            tt.fields.quit,
			}
			m.mining(tt.args.block, tt.args.difficulty)
		})
	}
}
