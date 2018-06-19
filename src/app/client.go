package app

type Client struct {
	ID []byte // i.e. Run-time port associated to container
	Peers [][]byte // Slice of IDs
	Chain []Block
}

// CLIENT-02, CLIENT-03, CLIENT-08, CLIENT-09
// should be implemented within this package

func (c *Client) Connect() error {
	// Try to correct
	return nil
}

func (c Client) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// CLIENT-04, CLIENT-06
	// To implement
	return nil
}

func (c Client) Disconnect() error {
	// Disconnect (RPC) to peers
	// CLIENT-05
	// To implement
	return nil
}
