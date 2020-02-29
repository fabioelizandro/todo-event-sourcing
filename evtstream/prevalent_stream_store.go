package evtstream

type PrevalentStreamStore interface {
	Load() ([]EventEnvelope, error)
	Write([]EventEnvelope) error
}
