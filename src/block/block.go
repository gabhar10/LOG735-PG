package block

import (
	"hash"
	"time"
)

const BLOCK_SIZE = 50

type Header struct {
	PreviousBlock hash.Hash64
	Nounce int64
	Date time.Time
}

type Block struct {
	Header Header
	Messages [BLOCK_SIZE]string
}