package evtstream

type prevalentEventStream struct {
	store       PrevalentStreamStore
	envelopes   []*prevalentEventEnvelope
	streamClock StreamClock
}

func NewPrevalentEventStream(store PrevalentStreamStore, envelopes []*prevalentEventEnvelope, streamClock StreamClock) *prevalentEventStream {
	return &prevalentEventStream{
		store:       store,
		envelopes:   envelopes,
		streamClock: streamClock,
	}
}

func (p *prevalentEventStream) FirstPosition() StreamPosition {
	return newPrevalentStreamPosition(0)
}

func (p *prevalentEventStream) Read(position StreamPosition) (EventEnvelope, error) {
	streamPosition := position.Value()

	count := int64(len(p.envelopes))

	if count == 0 {
		return nil, nil
	}

	if streamPosition+1 > count {
		return nil, nil
	}

	return p.envelopes[streamPosition], nil
}

func (p *prevalentEventStream) ReadByCorrelationID(correlationID string) ([]EventEnvelope, error) {
	correlatedEvents := make([]EventEnvelope, 0)
	for _, envelope := range p.envelopes {
		if envelope.Event().CorrelationID() == correlationID {
			correlatedEvents = append(correlatedEvents, envelope)
		}
	}

	return correlatedEvents, nil
}

func (p *prevalentEventStream) Write(events []Event) error {
	envelopes := []*prevalentEventEnvelope{}
	streamPosition := int64(len(p.envelopes))
	for _, event := range events {
		envelopes = append(envelopes, newPrevalentEventEnvelope(
			event,
			newPrevalentStreamPosition(streamPosition),
			p.streamClock.Now(),
		))
		streamPosition++
	}

	err := p.store.Write(envelopes)
	if err != nil {
		return err
	}

	p.envelopes = append(p.envelopes, envelopes...)

	return nil
}
