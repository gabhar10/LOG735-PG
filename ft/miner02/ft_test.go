package miner02

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"log"
	"reflect"
	"testing"
	"time"
)

// MINEUR-02: Lors d’une requête de connexion, le mineur doit fournir la chaîne de blocs.
func TestMiner02(t *testing.T) {
	const MinerID = "8889"
	const Client1ID = "8888"
	const Localhost = "127.0.0.1"
	const TestContent1 = "This is a test from Client1"

	t.Run("Lors d’une requête de connexion, le mineur doit fournir au client les 10 derniers blocs", func(t *testing.T) {
		// Create miner
		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: Client1ID,
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
		}

		// Create client1
		nodeChan1 := make(chan node.Message, 1)
		c1 := client.NewClient(Localhost, Client1ID, clientPeers, nil, nodeChan1).(*client.Client)
		err = c1.SetupRPC()
		if err != nil {
			t.Errorf("Error while setting up RPC: %v", err)
		}

		// Peer client1 with peers
		err = c1.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Make sure miner is peered!
		err = m.Peer()
		if err != nil {
			t.Errorf("Error while peering: %v", err)
		}

		// Let's make 10 blocks. 5 messages per block * 10 = 50 messages total
		for i := 0; i < node.BlockSize*10; i++ {
			msg := node.Message{
				Peer:    Client1ID,
				Content: TestContent1,
				Time:    time.Now().Format(time.RFC3339Nano),
			}

			err = c1.HandleUiMessage(msg)
			if err != nil {
				t.Errorf("Error while sending message: %v", err)
			}
		}

		for {
			c1Blocks := c1.GetBlocks()
			mBlocks := m.GetBlocks()
			if len(c1Blocks) == 10 && len(mBlocks) == 10 && reflect.DeepEqual(c1Blocks, mBlocks) {
				break
			}
			log.Printf("Client block size: %v\n", len(c1Blocks))
			time.Sleep(time.Second * 2)
		}
	})
}
