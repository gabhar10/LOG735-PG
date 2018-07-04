package node

import (
	"hash"
	"time"
)

type Header struct {
	PreviousBlock hash.Hash64
	Hash hash.Hash64
	Nounce int64
	Date time.Time
}

type Block struct {
	Header Header
	Messages [BlockSize]string
}