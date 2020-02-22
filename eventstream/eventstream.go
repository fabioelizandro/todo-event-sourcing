package eventstream

type EventStream interface {
	FirstPosition() StreamPosition
	Read(StreamPosition) (StreamReadResult, error)
	ReadByCorrelationID(string) ([]Event, error)
	Write([]Event) error
}

type StreamReadResult interface {
	Event() Event
	NextStreamPosition() StreamPosition
}

type StreamPosition interface {
	Before(StreamPosition) bool
	After(StreamPosition) bool
	Value() interface{}
}

type Event interface {
	Type() string
	CorrelationID() string
	Payload() ([]byte, error)
}
