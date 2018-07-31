package client09

import (
	"testing"
	"LOG735-PG/src/node"
	"LOG735-PG/src/client"
	"time"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/rpc"
)


// Un client doit afficher les messages dans l’ordre du contenue de chaque bloc
func TestClient09_parse(t *testing.T) {
	const ClientID = "8889"
	const MinerID = "8888"
	const Localhost = "127.0.0.1"
	t.Run("client receives and parses 2 blocks from miner", func(t *testing.T) {
		// Channel for communication


		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID,
			},
		}

		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)


		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID,
			},
		}

		uiChan := make(chan node.Message, node.BlockSize*2)
		c := client.NewClient(Localhost, ClientID, clientPeers, uiChan, nil, false).(*client.Client)

		err := m.SetupRPC()
		if err != nil {
			t.Errorf("Miner could not setup RPC: %v", err)
		}
		err = c.SetupRPC()
		if err != nil {
			t.Errorf("Client could not setup RPC: %v", err)
		}


		m.Start()

		err = c.Peer()
		if err != nil {
			t.Errorf("Error while peering the client to the miner: %v", err)
		}
		err = m.Peer()
		if err != nil {
			t.Errorf("Error while peering the miner to the client: %v", err)
		}

		m.ReceiveMessage("peer-1", "Test-1", "time-1", rpc.MessageType)
		m.ReceiveMessage("peer-2", "Test-2", "time-2", rpc.MessageType)
		m.ReceiveMessage("peer-3", "Test-3", "time-3", rpc.MessageType)
		m.ReceiveMessage("peer-4", "Test-4", "time-4", rpc.MessageType)
		m.ReceiveMessage("peer-5", "Test-5", "time-5", rpc.MessageType)

		received := false
		for i := 0; i < 5; i++ {
			if len(m.GetBlocks()) > 0 {
				received = true
				break
			}
			time.Sleep(time.Second * 2)
		}
		if !received {
			t.Errorf("Miner did not mine a block")
		}

		if len(uiChan) != 5 {
			t.Errorf("Channel should have 5 element in it; had %d", len(uiChan))
		}

		m.ReceiveMessage("peer-6", "Test-6", "time-6", rpc.MessageType)
		m.ReceiveMessage("peer-7", "Test-7", "time-7", rpc.MessageType)
		m.ReceiveMessage("peer-8", "Test-8", "time-8", rpc.MessageType)
		m.ReceiveMessage("peer-9", "Test-9", "time-9", rpc.MessageType)
		m.ReceiveMessage("peer-10", "Test-10", "time-10", rpc.MessageType)

		received = false
		for i := 0; i < 5; i++ {
			if len(m.GetBlocks()) > 1 {
				received = true
				break
			}
			time.Sleep(time.Second * 2)
		}
		if !received {
			t.Errorf("Miner did not mine a second block")
		}

		if len(uiChan) != 10 {
			t.Errorf("Channel should have 10 element in it; had %d", len(uiChan))
		}
	})
}
