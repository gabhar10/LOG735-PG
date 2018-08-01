package node

import (
	"net/rpc"
	"time"
)

// BlockSize is the size of a single block (e.g. amount of messages per block)
const BlockSize = 5

// MinBlocksReturnSize is the minimum amount of blocks to be returned by a miner
const MinBlocksReturnSize = 10

// MiningDifficulty is the amount of required zeroes heading a resulting hash
const MiningDifficulty = 2

// MessagesChannelSize is the size of the channel for incoming messages from other clients
const MessagesChannelSize = 100

// BlocksChannelSize is the size of the channel for incoming blocks from other miners
const BlocksChannelSize = 10

// MaliciousMinerSleepTime is the amount of time a malicious miner sleep before trying to find the nounce
const MaliciousMinerSleepTime = time.Microsecond * 10

// MinerSleepTime is the amount of time a miner sleep before trying to find the nounce
const MinerSleepTime = time.Second * 0

type Peer struct {
	Host string
	Port string
	Conn *rpc.Client
}
