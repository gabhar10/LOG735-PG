package client05

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"testing"
)


func TestClient05_connect(t *testing.T) {
	const MinerID = "8888"
	const ClientID = "8889"
	const TestContent = "This is a test"

	t.Run("Client connect to miner", func(t *testing.T) {
		// Create miner
		m := miner.NewMiner(MinerID, nil).(*miner.Miner)
		m.SetupRPC(MinerID)
		// Create client

		c := client.NewClient(ClientID, nil, nil, nil).(*client.Client)
		c.SetupRPC(ClientID)

		err := c.Connect(MinerID)
		t.Fatalf("Error returned while connecting to anchor: %v", err)



		found := false
		for i := 0; i < len(m.Peers); i++ {
			if m.Peers[i].Port == ClientID {
				if m.Peers[i].Conn == nil{
					t.Fatal("Miner has the peer but the connection is nil")
				} else {
					found = true
				}
			}
		}

		if found == false {
			t.Fatalf("Client does not have the peer")
		}

		for i := 0; i < len(c.Peers); i++ {
			if c.Peers[i].Port == ClientID {
				if c.Peers[i].Conn == nil{
					t.Fatal("Client has the peer but the connection is nil")
				} else {
					found = true
				}
			}
		}


	})
}

