package main

import (
	"LOG735-PG/src/client"
	"LOG735-PG/src/miner"
	"LOG735-PG/src/node"
	"fmt"
	"log"
	"os"
	"strings"
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
	var n node.Node
	role := os.Getenv("ROLE")

	peers := []*node.Peer{}
	for _, s := range strings.Split(os.Getenv("PEERS"), " ") {
		p := node.Peer{
			Host: fmt.Sprintf("node-%s", s),
			Port: s}

		peers = append(peers, &p)
	}

	switch role {
	case "client":
		n = client.NewClient(os.Getenv("PORT"), peers, nil, nil)
	case "miner":
		n = miner.NewMiner(os.Getenv("PORT"), peers)
	default:
		log.Fatalf("Unsupported role %s\n", role)
	}

	// Start RPC server
	err := n.SetupRPC()
	if err != nil {
		log.Fatal("RPC setup error:", err)
	}

	// bootstrap to each peer, blocking mechanism
	err = n.Peer()
	if err != nil {
		log.Fatal("Peering error:", err)
	}

	for {
		time.Sleep(time.Hour)
	}
}
