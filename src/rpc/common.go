package rpc

import "LOG735-PG/src/app"

type ConnectionRPC struct {
	PeerID []byte
}

type MessageRPC struct {
	ConnectionRPC
	Message string
}

type BlockRPC struct {
	ConnectionRPC
	app.Block
}