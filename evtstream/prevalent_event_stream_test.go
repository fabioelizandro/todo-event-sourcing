package evtstream_test

import (
	"fabioelizandro/todo-event-sourcing/evtstream"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_it_saves_events(t *testing.T) {
	clock := evtstream.NewFrozenClock(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	eventStream := evtstream.NewPrevalentEventStream(&FakePrevalentStore{}, nil, clock)

	assert.NoError(t, eventStream.Write([]evtstream.Event{
		&SomethingHappened{
			ID:   "1",
			Data: "Something",
		},
	}))

	envelope, err := eventStream.Read(eventStream.FirstPosition())
	assert.NoError(t, err)

	assert.Equal(t, envelope.Event(), &SomethingHappened{
		ID:   "1",
		Data: "Something",
	})
	assert.Equal(t, envelope.StreamPosition().Value(), int64(0))
	assert.Equal(t, envelope.Timestamp(), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
}

func Test_it_saves_events_to_store(t *testing.T) {
	clock := evtstream.NewFrozenClock(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	store := &FakePrevalentStore{}
	eventStream := evtstream.NewPrevalentEventStream(store, nil, clock)

	assert.NoError(t, eventStream.Write([]evtstream.Event{
		&SomethingHappened{
			ID:   "1",
			Data: "Something",
		},
	}))

	envelopes, err := store.Load()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(envelopes))
	assert.Equal(t, envelopes[0].Event(), &SomethingHappened{
		ID:   "1",
		Data: "Something",
	})
	assert.Equal(t, envelopes[0].StreamPosition().Value(), int64(0))
	assert.Equal(t, envelopes[0].Timestamp(), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
}

func Test_it_reads_by_correlation_id(t *testing.T) {
	clock := evtstream.NewFrozenClock(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	eventStream := evtstream.NewPrevalentEventStream(&FakePrevalentStore{}, nil, clock)

	assert.NoError(t, eventStream.Write([]evtstream.Event{
		&SomethingHappened{
			ID:   "1",
			Data: "Something",
		},
		&SomethingHappened{
			ID:   "2",
			Data: "Something",
		},
		&SomethingElseHappened{
			ID:   "1",
			Data: "Something Else",
		},
	}))

	envelopes, err := eventStream.ReadByCorrelationID("1")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(envelopes))
	assert.Equal(t, envelopes[0].Event(), &SomethingHappened{
		ID:   "1",
		Data: "Something",
	})
	assert.Equal(t, envelopes[1].Event(), &SomethingElseHappened{
		ID:   "1",
		Data: "Something Else",
	})
}

type FakePrevalentStore struct {
	memory []evtstream.EventEnvelope
	Err    error
}

func (m *FakePrevalentStore) Load() ([]evtstream.EventEnvelope, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.memory, nil
}

func (m *FakePrevalentStore) Write(envelopes []evtstream.EventEnvelope) error {
	if m.Err != nil {
		return m.Err
	}
	m.memory = append(m.memory, envelopes...)
	return nil
}
