package eventstream

// EventStream event stream interface
type EventStream interface {

	// ReadAggregate read all events from a specific aggregate id
	ReadAggregate(aggregateID string) ([]*EventEnvelope, error)

	// Write add events to the stream and ideally it would ensure uniqueness of AggregateID,AggregateType,AggregateVersion
	Write(events []*EventEnvelope) error
}

// EventEnvelope that's how a event looks like when we see it through the stream lens
type EventEnvelope struct {
	Type             string
	AggregateID      string
	AggregateType    string
	AggregateVersion int64
	Event            []byte
}

// InMemoryEventStream the in memory implementation of EventStream interface, use for test only purposes
type InMemoryEventStream struct {
	events []*EventEnvelope
}

func (stream *InMemoryEventStream) ReadAggregate(aggregateID string) ([]*EventEnvelope, error) {
	aggregateEvents := make([]*EventEnvelope, 0)
	for _, event := range stream.events {
		if event.AggregateID == aggregateID {
			aggregateEvents = append(aggregateEvents, event)
		}
	}

	return aggregateEvents, nil
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
