package evtstream

import "time"

type EventStream interface {
	FirstPosition() StreamPosition
	Read(StreamPosition) (EventEnvelope, error)
	ReadByCorrelationID(string) ([]EventEnvelope, error)
	Write([]Event) error
}

type StreamPosition interface {
	After(StreamPosition) bool
	Before(StreamPosition) bool
	Next() StreamPosition
	Value() int64
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
