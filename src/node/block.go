package node

import (
	"crypto/sha256"
	"time"
)

type Header struct {
<<<<<<< HEAD
	PreviousBlock [sha256.Size]byte
	Hash          [sha256.Size]byte
	Nounce        uint64
	Date          time.Time
=======
	PreviousBlock hash.Hash64
	Hash hash.Hash64
	Nounce int64
	Date time.Time
>>>>>>> cafaf6a91137f7faa5295e5bbcc83009b3c37049
}

type Block struct {
	Header   Header
	Messages [BlockSize]string
}
