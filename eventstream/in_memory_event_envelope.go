package eventstream

import (
	"time"
)

type inMemoryEventEnvelope struct {
	event          Event
	streamPosition StreamPosition
	timestamp      time.Time
}

func (i *inMemoryEventEnvelope) StreamPosition() StreamPosition {
	return i.streamPosition
}

func (i *inMemoryEventEnvelope) Timestamp() time.Time {
	return i.timestamp
}

func (i *inMemoryEventEnvelope) Event() Event {
	return i.event
}

func (i *inMemoryEventEnvelope) NextStreamPosition() StreamPosition {
	return &inMemoryStreamPosition{value: i.streamPosition.Value().(uint64) + 1}
}
