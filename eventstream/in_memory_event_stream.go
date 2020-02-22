package eventstream

func NewInMemoryEventStream() *inMemoryEventStream {
	return &inMemoryEventStream{}
}

type inMemoryEventStream struct {
	events []Event
}

func (s *inMemoryEventStream) FirstPosition() StreamPosition {
	return &inMemoryStreamPosition{value: 0}
}

func (s *inMemoryEventStream) Read(position StreamPosition) (StreamReadResult, error) {
	streamPosition := position.Value().(uint64)

	count := uint64(len(s.events))

	if count == 0 {
		return nil, nil
	}

	if streamPosition+1 > count {
		return nil, nil
	}

	return &inMemoryStreamReadResult{
		event:        s.events[streamPosition],
		nextPosition: &inMemoryStreamPosition{value: streamPosition + 1},
	}, nil
}

func (s *inMemoryEventStream) ReadByCorrelationID(correlationID string) ([]Event, error) {
	correlatedEvents := make([]Event, 0)
	for _, event := range s.events {
		if event.CorrelationID() == correlationID {
			correlatedEvents = append(correlatedEvents, event)
		}
	}

	return correlatedEvents, nil
}

func (s *inMemoryEventStream) Write(events []Event) error {
	s.events = append(s.events, events...)
	return nil
}

func (s *inMemoryEventStream) InMemoryReadAll() []Event {
	return s.events
}
