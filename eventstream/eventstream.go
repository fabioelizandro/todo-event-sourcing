package eventstream

type EventStream interface {
	Read(streamPosition uint64) (Event, error)
	ReadByCorrelationID(correlationID string) ([]Event, error)
	Write(events []Event) error
}

type Event interface {
	Type() string
	CorrelationID() string
	Payload() ([]byte, error)
}

type InMemoryEventStream struct {
	events []Event
}

func (stream *InMemoryEventStream) Read(streamPosition uint64) (Event, error) {
	count := uint64(len(stream.events))

	if count == 0 {
		return nil, nil
	}

	if streamPosition+1 > count {
		return nil, nil
	}

	return stream.events[streamPosition], nil
}

func (stream *InMemoryEventStream) ReadByCorrelationID(correlationID string) ([]Event, error) {
	correlatedEvents := make([]Event, 0)
	for _, event := range stream.events {
		if event.CorrelationID() == correlationID {
			correlatedEvents = append(correlatedEvents, event)
		}
	}

	return correlatedEvents, nil
}

func (stream *InMemoryEventStream) Write(events []Event) error {
	stream.events = append(stream.events, events...)
	return nil
}

func (stream *InMemoryEventStream) InMemoryReadAll() []Event {
	return stream.events
}

func NewInMemoryEventStream() *InMemoryEventStream {
	return &InMemoryEventStream{}
}
