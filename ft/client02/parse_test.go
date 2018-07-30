package client02

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/node"
	"testing"
	"LOG735-PG/src/miner"
	"time"
	"log"
)

//TODO : Should also check validity of Time of block
func TestClient02_parse(t *testing.T) {
	const MinerID = "8888"
	const ClientID = "8889"
	const TestContent = "This is a test"

	t.Run("Parse block from miner", func(t *testing.T) {

		minerPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: ClientID},
		}
		m := miner.NewMiner(MinerID, minerPeers).(*miner.Miner)

		// Create client
		clientPeers := []*node.Peer{
			&node.Peer{
				Host: "127.0.0.1",
				Port: MinerID},
		}
		// Channel for communication
		uiChan := make(chan node.Message, node.BlockSize)
		c := client.NewClient(ClientID, clientPeers, uiChan, nil).(*client.Client)

		//err := c.Peer()
		//if err != nil {
		//	t.Fatalf("Error while peering: %v", err)
		//}

		//var chain []node.Block
		var block1 node.Block
		var block2 node.Block


		for i := 0; i < node.BlockSize; i++{
			block1.Messages[i] = node.Message{"1","test-1",time.Now().Format(time.RFC3339Nano)}
		}


		log.Printf("%v", block1)
		go c.ReceiveBlock(block1)

		for i := 0; i < node.BlockSize; i++{
			msg := <-uiChan
			if msg.Content != "test-1"{
				t.Fatalf("Error message was %s; should be test-1", msg.Content)
			}
			if msg.Peer != "1"{
				t.Fatalf("Error peer was %s; should be 1", msg.Peer)
			}
		}



		for i := 0; i < node.BlockSize; i++{
			block2.Messages[i] = node.Message{"2","test-2",time.Now().Format(time.RFC3339Nano)}
		}

		go c.ReceiveBlock(block2)
		for i := 0; i < node.BlockSize; i++{
			msg := <-uiChan
			if msg.Content != "test-2"{
				t.Fatalf("Error message was %s; should be test-2", msg.Content)
			}
			if msg.Peer != "2"{
				t.Fatalf("Error peer was %s; should be 2", msg.Peer)
			}
		}






	})
}
