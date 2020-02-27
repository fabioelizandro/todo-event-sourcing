package eventstream

type PrevalentStreamStore interface {
	Load() ([]*prevalentEventEnvelope, error)
	Write([]*prevalentEventEnvelope) error
}
