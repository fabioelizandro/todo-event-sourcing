// +build integration

package evtstream_test

import (
	"fabioelizandro/todo-event-sourcing/evtstream"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_persists_events_to_disk(t *testing.T) {
	folder := fmt.Sprintf("/tmp/%s", uuid.New().String())
	store := evtstream.NewDiskPrevalentEventStore(folder, evtstream.NewInMemoryEventRegistry([]evtstream.Event{
		&SomethingHappened{},
		&SomethingElseHappened{},
	}))
	envelopes, err := store.Load()
	assert.NoError(t, err)

	stream := evtstream.NewPrevalentEventStream(store, envelopes)
	assert.NoError(t, stream.Write([]evtstream.Event{
		&SomethingHappened{ID: "1", Data: "foo"},
		&SomethingElseHappened{ID: "2", Data: "bar"},
	}))

	envelopes, err = store.Load()
	assert.NoError(t, err)

	stream = evtstream.NewPrevalentEventStream(store, envelopes)
	envelope1, err := stream.Read(stream.FirstPosition())
	assert.NoError(t, err)

	envelope2, err := stream.Read(envelope1.NextStreamPosition())
	assert.NoError(t, err)

	assert.Equal(t, []evtstream.Event{
		&SomethingHappened{ID: "1", Data: "foo"},
		&SomethingElseHappened{ID: "2", Data: "bar"},
	}, []evtstream.Event{
		envelope1.Event(),
		envelope2.Event(),
	})
}
