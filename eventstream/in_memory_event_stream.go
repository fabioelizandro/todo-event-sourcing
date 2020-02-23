package eventstream

import (
	"time"
)

func NewInMemoryEventStream() *inMemoryEventStream {
	return &inMemoryEventStream{}
}

type inMemoryEventStream struct {
	envelopes []EventEnvelope
}

func (s *inMemoryEventStream) FirstPosition() StreamPosition {
	return &inMemoryStreamPosition{value: 0}
}

func (s *inMemoryEventStream) Read(position StreamPosition) (EventEnvelope, error) {
	streamPosition := position.Value().(uint64)

	count := uint64(len(s.envelopes))

	if count == 0 {
		return nil, nil
	}

	if streamPosition+1 > count {
		return nil, nil
	}

	return s.envelopes[streamPosition], nil
}

func (s *inMemoryEventStream) ReadByCorrelationID(correlationID string) ([]EventEnvelope, error) {
	correlatedEvents := make([]EventEnvelope, 0)
	for _, envelope := range s.envelopes {
		if envelope.Event().CorrelationID() == correlationID {
			correlatedEvents = append(correlatedEvents, envelope)
		}
	}

	return correlatedEvents, nil
}

func (s *inMemoryEventStream) Write(events []Event) error {
	for _, event := range events {
		s.envelopes = append(s.envelopes, &inMemoryEventEnvelope{
			event:          event,
			streamPosition: &inMemoryStreamPosition{value: uint64(len(s.envelopes))},
			timestamp:      time.Now(),
		})
	}

	return nil
}

func (s *inMemoryEventStream) InMemoryReadAll() []Event {
	events := []Event{}

	for _, envelope := range s.envelopes {
		events = append(events, envelope.Event())
	}

	return events
}
