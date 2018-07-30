package miner03

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
	"time"
)

// MINEUR-03: Un mineur doit écouter pour les transactions des clients dans le réseau et pouvoir les accumuler dans une liste.
func TestMiner03(t *testing.T) {
	const MinerID = "8888"
	const ClientID = "8889"
	const TestContent = "This is a test"
	const Localhost = "localhost"

	t.Run("Send message to miner", func(t *testing.T) {
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
		err = c.SetupRPC()

		if err != nil {
			t.Errorf("Error while setting up RPC with client: %v", err)
		}

		err = c.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		err = m.Peer()
		if err != nil {
			t.Errorf("Error while miner peering: %v", err)
		}

		err = c.HandleUiMessage(node.Message{
			Peer:    ClientID,
			Content: TestContent,
			Time:    time.Now().Format(time.RFC3339Nano),
		})
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		if len(m.IncomingMsgChan) != 1 {
			t.Errorf("Message queue of miner should be 1")
		}

		msg := <-m.IncomingMsgChan

		timestamp, err := time.Parse(time.RFC3339Nano, msg.Time)
		if err != nil {
			t.Errorf("Error while parsing time: %v", err)
		}

		if msg.Content != TestContent || timestamp.After(time.Now()) {
			t.Errorf("Miner received wrong message")
		}
	})
}
