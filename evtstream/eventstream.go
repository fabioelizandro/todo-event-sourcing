package evtstream

import "time"

type EventStream interface {
	FirstPosition() StreamPosition
	Read(StreamPosition) (EventEnvelope, error)
	ReadByCorrelationID(string) ([]EventEnvelope, error)
	Write([]Event) error
}

type StreamPosition interface {
	Value() interface{}
}

type EventEnvelope interface {
	Event() Event
	StreamPosition() StreamPosition
	NextStreamPosition() StreamPosition
	Timestamp() time.Time
}

type Event interface {
	Type() string
	CorrelationID() string
}
