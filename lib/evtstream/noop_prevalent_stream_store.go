package evtstream

type noopPrevalentStreamStore struct {
}

func newNoopPrevalentStreamStore() *noopPrevalentStreamStore {
	return &noopPrevalentStreamStore{}
}

func (n *noopPrevalentStreamStore) Load() ([]EventEnvelope, error) {
	return []EventEnvelope{}, nil
}

func (n *noopPrevalentStreamStore) Write([]EventEnvelope) error {
	return nil
}
