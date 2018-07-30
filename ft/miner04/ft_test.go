package miner04

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"testing"
	"time"
)

// MINEUR-04: Un mineur doit respecter l’ordre chronologique de l’envoi des messages par les clients lors de la création de blocs.
func TestMiner04(t *testing.T) {
	const MinerID = "8889"
	const Client1ID = "8888"
	const Client2ID = "8887"
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
		err := m.SetupRPC()
		if err != nil {
			t.Errorf("Error while trying to setup RPC: %v", err)
		}

		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}

		// Create client1
		nodeChan1 := make(chan node.Message, 1)
		c1 := client.NewClient(Client1ID, clientPeers, nil, nodeChan1).(*client.Client)

		err = c1.SetupRPC()
		if err != nil {
			t.Errorf("Error while setting up RPC with client: %v", err)
		}

		err = c1.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Create client2
		nodeChan2 := make(chan node.Message, 1)
		c2 := client.NewClient(Client2ID, clientPeers, nil, nodeChan2).(*client.Client)

		err = c2.SetupRPC()
		if err != nil {
			t.Errorf("Error while setting up RPC with client: %v", err)
		}

		err = c2.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Peer miner
		err = m.Peer()
		if err != nil {
			t.Errorf("Error during miner peering: %v", err)
		}

		msg := node.Message{
			Peer:    Client1ID,
			Content: TestContent1,
			Time:    time.Now().Format(time.RFC3339Nano),
		}

		err = c1.HandleUiMessage(msg)
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent3,
			Time:    time.Now().Format(time.RFC3339Nano),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client1ID,
			Content: TestContent2,
			Time:    time.Now().Format(time.RFC3339Nano),
		}

		err = c1.HandleUiMessage(msg)
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent4,
			Time:    time.Now().Format(time.RFC3339Nano),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		if len(m.IncomingMsgChan) != 4 {
			t.Errorf("Message queue of miner should be 4")
		}

		msg = <-m.IncomingMsgChan
		timestamp, err := time.Parse(time.RFC3339Nano, msg.Time)
		if err != nil {
			t.Errorf("Error while parsing time: %v", err)
		}
		if msg.Content != TestContent1 || timestamp.After(time.Now()) {
			t.Errorf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		timestamp, err = time.Parse(time.RFC3339Nano, msg.Time)
		if err != nil {
			t.Errorf("Error while parsing time: %v", err)
		}
		if msg.Content != TestContent3 || timestamp.After(time.Now()) {
			t.Errorf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		timestamp, err = time.Parse(time.RFC3339Nano, msg.Time)
		if err != nil {
			t.Errorf("Error while parsing time: %v", err)
		}
		if msg.Content != TestContent2 || timestamp.After(time.Now()) {
			t.Errorf("Miner received wrong message")
		}

		msg = <-m.IncomingMsgChan
		timestamp, err = time.Parse(time.RFC3339Nano, msg.Time)
		if err != nil {
			t.Errorf("Error while parsing time: %v", err)
		}
		if msg.Content != TestContent4 || timestamp.After(time.Now()) {
			t.Errorf("Miner received wrong message")
		}
	})
}
