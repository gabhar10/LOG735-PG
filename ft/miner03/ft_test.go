package miner03

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
	"time"
)

func TestMiner03(t *testing.T) {
	const MinerID = "8888"
	const ClientID = "8889"
	const TestContent = "This is a test"

	t.Run("Send message to miner", func(t *testing.T) {
		// Create miner
		minerPeers := []node.Peer{
			node.Peer{
				Host: "127.0.0.1",
				Port: ClientID},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)
		m.SetupRPC(MinerID)
		// Create client
		clientPeers := []node.Peer{
			node.Peer{
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

		nodeChan <- node.Message{
			Peer:    ClientID,
			Content: TestContent,
			Time:    time.Now(),
		}

		err = c.HandleUiMessage()
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		if len(m.IncomingMsgChan) != 1 {
			t.Fatalf("Message queue of miner should be 1")
		}

		msg := <-m.IncomingMsgChan
		if msg.Content != TestContent || msg.Time.After(time.Now()) {
			t.Fatalf("Miner received wrong message")
		}
	})
}
