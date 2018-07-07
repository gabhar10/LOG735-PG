package app

import (
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"testing"
	"time"
)

func TestMiner_mining(t *testing.T) {
	type fields struct {
		ID                string
		blocks            []node.Block
		peers             []string
		rpcHandler        *brpc.NodeRPC
		incomingMsgChan   chan Message
		incomingBlockChan chan node.Block
		quit              chan bool
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
				make(chan Message, 100),
				make(chan node.Block, 10),
				make(chan bool),
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
			}
			m.mining()
			if m.blocks[1].Header.Date.Equal(time.Time{}) {
				t.Errorf("date is not defined")
			}
		})
	}
}
