package client02

import (
	"testing"
)

//TODO : Should also check validity of Time of block
func TestClient02_parse(t *testing.T) {
	const ClientID = "8889"

	t.Run("Parse block", func(t *testing.T) {
		/*// Channel for communication
		uiChan := make(chan node.Message, node.BlockSize)
		c := client.NewClient(ClientID, nil, uiChan, nil).(*client.Client)

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
		}*/
	})
}
