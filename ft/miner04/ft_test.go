package miner04

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
	"time"
)

func TestMiner04(t *testing.T) {
	const MinerID = "8889"
	const Client1ID = "8888"
	const Client2ID = "8889"
	const TestContent1 = "This is a test from Client1"
	const TestContent2 = "This is a another test from Client1"
	const TestContent3 = "This is a test from Client2"
	const TestContent4 = "This is a another test from Client2"

	t.Run("Send message to miner", func(t *testing.T) {
		// Create miner
		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: Client1ID,
			},
			&node.Peer{
				Host: "127.0.0.1",
				Port: Client2ID,
			},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)
		m.SetupRPC(MinerID)

		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}

		// Create client1
		nodeChan1 := make(chan node.Message, 1)
		c1 := client.NewClient(Client1ID, clientPeers, nil, nodeChan1).(*client.Client)
		err := c1.Peer()
		if err != nil {
			t.Fatalf("Error while peering: %v", err)
		}

		// Create client2
		nodeChan2 := make(chan node.Message, 1)
		c2 := client.NewClient(Client2ID, clientPeers, nil, nodeChan2).(*client.Client)
		err = c2.Peer()
		if err != nil {
			t.Fatalf("Error while peering: %v", err)
		}

		nodeChan1 <- node.Message{
			Peer:    Client1ID,
			Content: TestContent1,
			Time:    time.Now(),
		}

		err = c1.HandleUiMessage()
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		nodeChan2 <- node.Message{
			Peer:    Client1ID,
			Content: TestContent3,
			Time:    time.Now(),
		}

		err = c2.HandleUiMessage()
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		nodeChan1 <- node.Message{
			Peer:    Client1ID,
			Content: TestContent2,
			Time:    time.Now(),
		}

		err = c1.HandleUiMessage()
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		nodeChan2 <- node.Message{
			Peer:    Client1ID,
			Content: TestContent4,
			Time:    time.Now(),
		}

		err = c2.HandleUiMessage()
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		if len(m.IncomingMsgChan) != 4 {
			t.Fatalf("Message queue of miner should be 4")
		}

		msg := <-m.IncomingMsgChan
		if msg.Content != TestContent1 || msg.Time.After(time.Now()) {
			t.Fatalf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		if msg.Content != TestContent3 || msg.Time.After(time.Now()) {
			t.Fatalf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		if msg.Content != TestContent2 || msg.Time.After(time.Now()) {
			t.Fatalf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		if msg.Content != TestContent4 || msg.Time.After(time.Now()) {
			t.Fatalf("Miner received wrong message")
		}
	})
}