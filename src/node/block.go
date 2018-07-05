package node

import (
	"crypto/sha256"
	"time"
)

type Header struct {
	PreviousBlock [sha256.Size]byte
	Hash          [sha256.Size]byte
	Nounce        uint64
	Date          time.Time
}

type Block struct {
	Header   Header
	Messages [BlockSize]string
}
