package client05

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
)

func TestClient05_disconnect(t *testing.T) {
	const MinerID = "8888"
	const ClientID = "8889"

	t.Run("Send message to miner", func(t *testing.T) {
		// Create miner
		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)
		m.SetupRPC(MinerID)
		// Create client
		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}
		// Channel for communication
		nodeChan := make(chan node.Message, 1)
		c := client.NewClient(ClientID, clientPeers, nil, nodeChan).(*client.Client)
		err := c.Peer()
		if err != nil {
			t.Fatalf("Error while peering: %v", err)
		}

		if len(m.Peers) == 0 {
			t.Fatalf("Miner's peers slice is empty")
		}

		if len(c.Peers) == 0 {
			t.Fatalf("Miner's peers slice is empty")
		}

		err = c.Disconnect()
		if err != nil {
			t.Fatalf("Error returned while disconnecting: %v", err)
		}

		if len(m.Peers) > 0 {
			t.Fatalf("Miner still has the connection of the peer")
		}

		if len(c.Peers) > 0 {
			t.Fatal("Client still has the connection of the peer")
		}
	})
}
