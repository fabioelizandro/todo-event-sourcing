package eventstream

// EventStream event stream interface
type EventStream interface {

	// ReadAggregate read all events from a specific aggregate id
	ReadAggregate(aggregateID string, eventHandler eventHandler) error

	// Write add events to the stream and ideally it would ensure uniqueness of AggregateID,AggregateType,AggregateVersion
	Write(events []*EventEnvelope) error
}

// EventEnvelope that's how a event looks like when we see it through the stream lens
type EventEnvelope struct {
	Type             string
	AggregateID      string
	AggregateType    string
	AggregateVersion int64
	Event            Event
}

type Event interface {
	Payload() ([]byte, error)
}

// eventHandler is a func to handle events when reading from the event stream
type eventHandler = func(event *EventEnvelope)

// InMemoryEventStream the in memory implementation of EventStream interface, use for test only purposes
type InMemoryEventStream struct {
	events []*EventEnvelope
}

func (stream *InMemoryEventStream) ReadAggregate(aggregateID string, handler eventHandler) error {
	for _, event := range stream.events {
		if event.AggregateID == aggregateID {
			handler(event)
		}
	}

	return nil
}

func (stream *InMemoryEventStream) Write(events []*EventEnvelope) error {
	// TODO: ensure uniqueness of events (aggregateID + aggregateVersion)
	stream.events = append(stream.events, events...)
	return nil
}

func (stream *InMemoryEventStream) InMemoryReadAll() []*EventEnvelope {
	return stream.events
}

func NewInMemoryEventStream() *InMemoryEventStream {
	return &InMemoryEventStream{}
}
