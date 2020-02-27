package evtstream

import (
	"time"
)

type prevalentEventEnvelope struct {
	event          Event
	streamPosition StreamPosition
	timestamp      time.Time
}

func newPrevalentEventEnvelope(event Event, streamPosition StreamPosition, timestamp time.Time) *prevalentEventEnvelope {
	return &prevalentEventEnvelope{
		event:          event,
		streamPosition: streamPosition,
		timestamp:      timestamp,
	}
}

func (i *prevalentEventEnvelope) StreamPosition() StreamPosition {
	return i.streamPosition
}

func (i *prevalentEventEnvelope) Timestamp() time.Time {
	return i.timestamp
}

func (i *prevalentEventEnvelope) Event() Event {
	return i.event
}

func (i *prevalentEventEnvelope) NextStreamPosition() StreamPosition {
	return i.streamPosition.Next()
}
