package main

import (
	"LOG735-PG/src/app"
	"LOG735-PG/src/node"
	"log"
	"os"
	"time"
)

func main() {
	// Check for empty environment variables
	for _, role := range []string{"ROLE", "PORT", "PEERS"} {
		env := os.Getenv(role)
		if env == "" {
			log.Fatalf("Environment variable %s is empty\n", role)
		}
	}

	// Define RPC handler
	var node node.Node
	role := os.Getenv("ROLE")
	switch role {
	case "client":
		node = app.NewClient(os.Getenv("PORT"), os.Getenv("PEERS"), nil)
	case "miner":
		node = app.NewMiner(os.Getenv("PORT"), os.Getenv("PEERS"))
	default:
		log.Fatalf("Unsupported role %s\n", role)
	}

	// Start RPC server
	nodePort := os.Getenv("PORT")
	err := node.SetupRPC(nodePort)
	if err != nil {
		log.Fatal("RPC setup error:", err)
	}

	// bootstrap to each peer, blocking mechanism
	err = node.Peer()
	if err != nil {
		log.Fatal("Peering error:", err)
	}

	// fmt.Println("Hello")

	// miner := app.Miner{}
	// miner.CreateBlock()

	for {
		time.Sleep(time.Hour)
	}
}
