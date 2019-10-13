package eventstream

type EventStream interface {
	Read(eventID uint64) (Event, error)
	ReadAggregate(aggregateID string) ([]Event, error)
	Write(events []Event) error
}

type Event interface {
	Type() string
	AggregateID() string
	AggregateType() string
	Payload() ([]byte, error)
}

type InMemoryEventStream struct {
	events []Event
}

func (stream *InMemoryEventStream) Read(eventID uint64) (Event, error) {
	count := uint64(len(stream.events))

	if count == 0 {
		return nil, nil
	}

	if eventID+1 > count {
		return nil, nil
	}

	return stream.events[eventID], nil
}

func (stream *InMemoryEventStream) ReadAggregate(aggregateID string) ([]Event, error) {
	aggregateEvents := make([]Event, 0)
	for _, event := range stream.events {
		if event.AggregateID() == aggregateID {
			aggregateEvents = append(aggregateEvents, event)
		}
	}

	return aggregateEvents, nil
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
