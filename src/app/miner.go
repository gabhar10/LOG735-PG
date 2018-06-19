package app

type Miner struct {
	Chain []Block // MINEUR-07
}

// MINEUR-12

func (m Miner) Broadcast() error {
	// DeliverMessage (RPC) to peers
	// MINEUR-04
	// To implement
	return nil
}

func (m *Miner) CreateBlock() error {
	// MINEUR-10
	// MINEUR-14
	// To implement
	header := &Header{}
	err := m.findNounce(header, uint64(0))
	if err != nil {
		return err
	}

	// Broadcast to all peers
	// MINEUR-06
	return nil
}

func (m Miner) findNounce(header *Header, difficulty uint64) error {
	// MINEUR-05
	return nil
}