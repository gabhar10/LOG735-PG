package client

import (
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"testing"
)

func TestClient_Peer(t *testing.T) {
	type fields struct {
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
				ID: "9001",
				peers: func() []*node.Peer {
					driver := node.Peer{Host: "127.0.0.1", Port: "9001"}
					m := miner.NewMiner("9002", []*node.Peer{&driver}).(*miner.Miner)
					m.SetupRPC("9002")
					return []*node.Peer{&node.Peer{Host: "127.0.0.1", Port: "9002"}}
				}(),
				uiChannel:   nil,
				nodeChannel: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.ID, tt.fields.peers, tt.fields.uiChannel, tt.fields.nodeChannel)
			if err := c.Peer(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Peer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
