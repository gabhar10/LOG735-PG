package client02

import (
	"testing"
	"LOG735-PG/src/node"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/client"
	"time"
)

//TODO : Should also check validity of Time of block
func TestClient02_getBlock(t *testing.T) {
	const ClientID1 = "8889"
	const ClientID2 = "8890"
	const ClientID3 = "8891"
	const MinerID = "9000"
	const Localhost = "127.0.0.1"

	const TestContent1 = "Test-1"
	const TestContent2 = "Test-2"
	const TestContent3 = "Test-3"
	const TestContent4 = "Test-4"
	const TestContent5 = "Test-5"

	t.Run("Broadcast of messages to 2 client from 1 miner and getblock for a third client", func(t *testing.T) {

		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID1,
			},
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID2,
			},
		}

		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID,
			},
		}

		uiChan1 := make(chan node.Message, node.BlockSize)
		uiChan2 := make(chan node.Message, node.BlockSize)
		c1 := client.NewClient(Localhost, ClientID1, clientPeers, uiChan1, nil).(*client.Client)
		c2 := client.NewClient(Localhost, ClientID2, clientPeers, uiChan2, nil).(*client.Client)
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)



		err := m.SetupRPC()
		if err != nil {
			t.Errorf("Miner could not setup RPC: %v", err)
		}
		m.Start()
		err = c1.SetupRPC()
		if err != nil {
			t.Errorf("Client-1 could not setup RPC: %v", err)
		}
		err = c2.SetupRPC()
		if err != nil {
			t.Errorf("Client-2 could not setup RPC: %v", err)
		}
		err = c1.Peer()
		if err != nil {
			t.Errorf("Error while peering the first client to the miner: %v", err)
		}
		err = c2.Peer()
		if err != nil {
			t.Errorf("Error while peering the second client to the miner: %v", err)
		}
		err = m.Peer()
		if err != nil {
			t.Errorf("Error while peering the miner to the client: %v", err)
		}

		c1.HandleUiMessage( node.Message{Content:TestContent1, Peer:c1.ID, Time:time.Now().Format(time.RFC3339Nano)})
		c2.HandleUiMessage( node.Message{Content:TestContent2, Peer:c2.ID, Time:time.Now().Format(time.RFC3339Nano)})
		c2.HandleUiMessage( node.Message{Content:TestContent3, Peer:c2.ID, Time:time.Now().Format(time.RFC3339Nano)})
		c1.HandleUiMessage( node.Message{Content:TestContent4, Peer:c1.ID, Time:time.Now().Format(time.RFC3339Nano)})
		c2.HandleUiMessage( node.Message{Content:TestContent5, Peer:c2.ID, Time:time.Now().Format(time.RFC3339Nano)})

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

		if len(uiChan1) != 5 {
			t.Errorf("Channel of client 1 should have 5 element in it; had %d", len(uiChan1))
		}
		if len(uiChan2) != 5 {
			t.Errorf("Channel of client 2 should have 5 element in it; had %d", len(uiChan2))
		}

		// Check message client 1 received
		msg := <- uiChan1
		if msg.Content != TestContent1 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent1,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan1
		if msg.Content != TestContent2 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent2,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan1
		if msg.Content != TestContent3 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent3,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan1
		if msg.Content != TestContent4 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent4,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan1
		if msg.Content != TestContent5|| msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent5,
				ClientID2,
				msg.Content,
				msg.Peer)
		}

		// Check message client 2 received
		msg = <- uiChan2
		if msg.Content != TestContent1 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent1,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan2
		if msg.Content != TestContent2 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent2,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan2
		if msg.Content != TestContent3 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent3,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan2
		if msg.Content != TestContent4 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent4,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan2
		if msg.Content != TestContent5 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent5,
				ClientID2,
				msg.Content,
				msg.Peer)
		}

		uiChan3 := make(chan node.Message, node.BlockSize)
		c3 := client.NewClient(Localhost, ClientID3, clientPeers, uiChan3, nil).(*client.Client)
		c3.Connect(Localhost, MinerID)

		msg = <- uiChan3
		if msg.Content != TestContent1 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent1,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan3
		if msg.Content != TestContent2 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent2,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan3
		if msg.Content != TestContent3 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent3,
				ClientID2,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan3
		if msg.Content != TestContent4 || msg.Peer != c1.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent4,
				ClientID1,
				msg.Content,
				msg.Peer)
		}
		msg = <- uiChan3
		if msg.Content != TestContent5 || msg.Peer != c2.ID {
			t.Errorf("Message not in order; should be %s from %s; Was %s from %s",
				TestContent5,
				ClientID2,
				msg.Content,
				msg.Peer)
		}

	})
}
