package miner06

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"reflect"
	"testing"
	"time"
)

// MINEUR-06: Après avoir découvert une solution valide à la création d’un bloc, un mineur doit diffuser ce bloc à ses voisins.
func TestMiner06(t *testing.T) {
	const MinerID = "8889"
	const Client1ID = "8888"
	const Client2ID = "8887"
	const MinerPeerID = "8886"
	const Localhost = "localhost"
	const TestContent1 = "This is a test from Client1"
	const TestContent2 = "This is a another test from Client1"
	const TestContent3 = "This is a test from Client2"
	const TestContent4 = "This is a another test from Client2"
	const TestContent5 = "This is again another test from Client2"

	t.Run("Après avoir découvert une solution valide à la création d’un bloc, un mineur doit diffuser ce bloc à ses voisins.", func(t *testing.T) {
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
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerPeerID,
			},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)
		err := m.SetupRPC()
		if err != nil {
			t.Errorf("Miner could not setup RPC: %v", err)
		}
		m.Start()

		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID,
			},
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerPeerID,
			},
		}

		// Create client1
		nodeChan1 := make(chan node.Message, 1)
		c1 := client.NewClient(Localhost, Client1ID, clientPeers, nil, nodeChan1, false).(*client.Client)
		err = c1.SetupRPC()
		if err != nil {
			t.Errorf("Error while setting up RPC: %v", err)
		}

		// Create miner peer
		mPeerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID,
			},
			&node.Peer{
				Host: "127.0.0.1",
				Port: Client1ID,
			},
		}

		mPeer := miner.NewMiner(MinerPeerID, mPeerPeers)
		err = mPeer.SetupRPC()
		if err != nil {
			t.Errorf("Error while setting up RPC: %v", err)
		}
		err = mPeer.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Peer client1 with both miner driver and miner peer
		err = c1.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Create client2
		nodeChan2 := make(chan node.Message, 1)
		c2 := client.NewClient(Localhost, Client2ID, clientPeers, nil, nodeChan2, false).(*client.Client)
		c2.SetupRPC()
		err = c2.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Make sure miner is peered!
		err = m.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
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

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent5,
			Time:    time.Now().Format(time.RFC3339Nano),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Errorf("Error while sending message: %v", err)
		}

		// Wait maximum 25 seconds for miner to mine its block

		time.Sleep(time.Second * 5)

		for {
			if len(c1.GetBlocks()) > 1 || len(c2.GetBlocks()) > 1 || len(mPeer.GetBlocks()) > 1 {
				t.Errorf("All nodes should have blocks of size 1")
			}
			if len(c1.GetBlocks()) == 1 && len(c2.GetBlocks()) == 1 && len(mPeer.GetBlocks()) == 1 {
				break
			}
			time.Sleep(time.Second * 2)
		}

		c1Blocks := c1.GetBlocks()
		c2Blocks := c2.GetBlocks()
		mPeerBlocks := mPeer.GetBlocks()

		if !reflect.DeepEqual(c1Blocks, c2Blocks) {
			t.Errorf("Both clients did not receive the same block")
		}

		if !reflect.DeepEqual(c2Blocks, mPeerBlocks) {
			t.Errorf("Miner and clients did not receive the same block: %v != %v", c2Blocks, mPeerBlocks)
		}

		if c1Blocks[0].Messages[0].Content != TestContent1 {
			t.Errorf("First message is not in right order")
		}
		if c1Blocks[0].Messages[1].Content != TestContent3 {
			t.Errorf("Second message is not in right order")
		}
		if c1Blocks[0].Messages[2].Content != TestContent2 {
			t.Errorf("Third message is not in right order")
		}
		if c1Blocks[0].Messages[3].Content != TestContent4 {
			t.Errorf("Fourth message is not in right order")
		}
		if c1Blocks[0].Messages[4].Content != TestContent5 {
			t.Errorf("Fifth message is not in right order")
		}
	})
}
