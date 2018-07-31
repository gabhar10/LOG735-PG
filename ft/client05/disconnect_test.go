package client05

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
)

func TestClient05_disconnect(t *testing.T) {
	const MinerID = "9888"
	const ClientID = "9889"
	const Localhost = "localhost"
	// CLIENT-05: L’application client doit permettre à l’utilisateur de se déconnecter


	t.Run("L’application client doit permettre à l’utilisateur de se déconnecter", func(t *testing.T) {
		// Create miner
		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)
		err := m.SetupRPC()
		if err != nil {
			t.Errorf("Error while trying to setup RPC: %v", err)
		}
		// Create client
		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}
		// Channel for communication
		nodeChan := make(chan node.Message, 1)

		c := client.NewClient(Localhost, ClientID, clientPeers, nil, nodeChan).(*client.Client)
		err = c.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		if len(m.Peers) == 0 {
			t.Errorf("Miner's peers slice is empty")
		}

		if len(c.Peers) == 0 {
			t.Errorf("Miner's peers slice is empty")
		}

		err = c.Disconnect()

		if err != nil {
			t.Errorf("Error returned while disconnecting: %v", err)
		}

		if len(m.Peers) > 0 {
			t.Errorf("Miner still has the connection of the peer")
		}

		if len(c.Peers) > 0 {
			t.Errorf("Client still has the connection of the peer")
		}
	})
}
