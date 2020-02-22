package eventstream

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type prevalentEventStream struct {
	sync.Mutex
	memory         *inMemoryEventStream
	streamPosition uint64
	streamFolder   string
}

type streamBatch struct {
	Events []*streamBatchEvent
}

type streamBatchEvent struct {
	Position uint64
	Type     string
	Event    []byte
}

// TODO: break this down into multiple functions or change it to be lazy
func LoadPrevalentEventStream(streamFolder string, eventRegistryMap map[string]Event) (*prevalentEventStream, error) {
	files, err := ioutil.ReadDir(streamFolder)
	if err != nil {
		return nil, err
	}

	allStreamBatchEvents := []*streamBatchEvent{}

	for _, f := range files { // TODO: read all files in parallel
		content, err := ioutil.ReadFile(f.Name())
		if err != nil {
			return nil, err
		}

		streamBatchEvents := []*streamBatchEvent{}
		err = json.Unmarshal(content, streamBatchEvents)
		if err != nil {
			return nil, err
		}

		allStreamBatchEvents = append(allStreamBatchEvents, streamBatchEvents...)
	}

	sort.Slice(allStreamBatchEvents, func(i, j int) bool {
		return allStreamBatchEvents[i].Position < allStreamBatchEvents[j].Position
	})

	events := []Event{}
	for _, streamEvent := range allStreamBatchEvents {
		t, ok := eventRegistryMap[streamEvent.Type]
		if !ok {
			return nil, fmt.Errorf("event type not mapped")
		}

		tCopy := reflect.New(reflect.ValueOf(t).Elem().Type()).Interface().(Event)

		err := json.Unmarshal(streamEvent.Event, tCopy)
		if err != nil {
			return nil, err
		}

		events = append(events, tCopy)
	}

	memory := NewInMemoryEventStream()
	err = memory.Write(events)
	if err != nil {
		return nil, err
	}

	return &prevalentEventStream{
		memory:         memory,
		streamPosition: uint64(len(events) - 1),
		streamFolder:   streamFolder,
	}, nil
}

func (p *prevalentEventStream) FirstPosition() StreamPosition {
	return p.memory.FirstPosition()
}

func (p *prevalentEventStream) Read(position StreamPosition) (StreamReadResult, error) {
	return p.memory.Read(position)
}

func (p *prevalentEventStream) ReadByCorrelationID(id string) ([]Event, error) {
	return p.memory.ReadByCorrelationID(id)
}

func (p *prevalentEventStream) Write(events []Event) error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	err := p.saveToDist(events)
	if err != nil {
		return err
	}

	return p.memory.Write(events)
}

func (p *prevalentEventStream) saveToDist(events []Event) error {
	streamPosition := p.streamPosition

	streamBatchEvents := []*streamBatchEvent{}
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		streamBatchEvents = append(streamBatchEvents, &streamBatchEvent{
			Position: streamPosition,
			Type:     event.Type(),
			Event:    payload,
		})

		streamPosition++
	}

	batchPayload, err := json.Marshal(streamBatch{Events: streamBatchEvents})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		fmt.Sprintf("%s/%s", p.streamFolder, uuid.New().String()),
		batchPayload,
		0644,
	)
	if err != nil {
		return err
	}

	p.streamPosition = streamPosition
	return nil
}

func loadFromDisk(streamFolder string) ([]Event, error) {

}
