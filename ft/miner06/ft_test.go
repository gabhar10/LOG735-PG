package miner06

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"reflect"
	"testing"
	"time"
)

func TestMiner06(t *testing.T) {
	const MinerID = "8889"
	const Client1ID = "8888"
	const Client2ID = "8887"
	const TestContent1 = "This is a test from Client1"
	const TestContent2 = "This is a another test from Client1"
	const TestContent3 = "This is a test from Client2"
	const TestContent4 = "This is a another test from Client2"
	const TestContent5 = "This is again another test from Client2"

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
		err := m.SetupRPC(MinerID)
		if err != nil {
			t.Fatalf("Miner could not setup RPC: %v", err)
		}
		m.Start()

		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}

		// Create client1
		nodeChan1 := make(chan node.Message, 1)
		c1 := client.NewClient(Client1ID, clientPeers, nil, nodeChan1).(*client.Client)
		err = c1.Peer()
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

		msg := node.Message{
			Peer:    Client1ID,
			Content: TestContent1,
			Time:    time.Now(),
		}

		err = c1.HandleUiMessage(msg)
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent3,
			Time:    time.Now(),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client1ID,
			Content: TestContent2,
			Time:    time.Now(),
		}

		err = c1.HandleUiMessage(msg)
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent4,
			Time:    time.Now(),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		msg = node.Message{
			Peer:    Client2ID,
			Content: TestContent5,
			Time:    time.Now(),
		}

		err = c2.HandleUiMessage(msg)
		if err != nil {
			t.Fatalf("Error while sending message: %v", err)
		}

		// Wait 25 seconds for miner to mine its block
		received := false
		for i := 0; i < 5; i++ {
			if len(c1.GetBlocks()) > 0 && len(c2.GetBlocks()) > 0 {
				received = true
				break
			}
			time.Sleep(time.Second * 2)
		}

		if !received {
			t.Fatalf("Never received a block")
		}

		c1Blocks := c1.GetBlocks()
		c2Blocks := c2.GetBlocks()

		if !reflect.DeepEqual(c1Blocks, c2Blocks) {
			t.Fatalf("Both clients did not receive the same block")
		}

		if c1Blocks[0].Messages[0].Content != TestContent1 {
			t.Fatalf("First message is not in right order")
		}
		if c1Blocks[0].Messages[1].Content != TestContent3 {
			t.Fatalf("Second message is not in right order")
		}
		if c1Blocks[0].Messages[2].Content != TestContent2 {
			t.Fatalf("Third message is not in right order")
		}
		if c1Blocks[0].Messages[3].Content != TestContent4 {
			t.Fatalf("Fourth message is not in right order")
		}
		if c1Blocks[0].Messages[4].Content != TestContent5 {
			t.Fatalf("Fifth message is not in right order")
		}
	})
}
