package eventstream

type noopPrevalentStreamStore struct {
}

func newNoopPrevalentStreamStore() *noopPrevalentStreamStore {
	return &noopPrevalentStreamStore{}
}

func (n *noopPrevalentStreamStore) Load() ([]*prevalentEventEnvelope, error) {
	return []*prevalentEventEnvelope{}, nil
}

func (n *noopPrevalentStreamStore) Write([]*prevalentEventEnvelope) error {
	return nil
}
