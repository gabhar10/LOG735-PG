package main

import (
	"fmt"
	"os"
)

func main() {
	role := os.Getenv("ROLE")
	peers := os.Getenv("PEERS")

	fmt.Printf("Role: %s\nPeers:%s\n", role, peers)
}