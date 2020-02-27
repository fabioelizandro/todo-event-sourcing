package evtstream

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
)

type diskPrevalentStreamStore struct {
	folder   string
	registry EventRegistry
}

func NewDiskPrevalentEventStore(folder string, registry EventRegistry) *diskPrevalentStreamStore {
	return &diskPrevalentStreamStore{
		folder:   folder,
		registry: registry,
	}
}

func (d *diskPrevalentStreamStore) Load() ([]*prevalentEventEnvelope, error) {
	err := os.MkdirAll(d.folder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(d.folder)
	if err != nil {
		return nil, err
	}

	envelopes := []*prevalentEventEnvelope{}
	for _, f := range files {
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", d.folder, f.Name()))
		if err != nil {
			return nil, err
		}

		decodedEnvelopes, err := d.decode(content)
		if err != nil {
			return nil, err
		}

		envelopes = append(envelopes, decodedEnvelopes...)
	}

	sort.Slice(envelopes, func(i, j int) bool {
		return envelopes[i].StreamPosition().Value().(uint64) < envelopes[j].StreamPosition().Value().(uint64)
	})

	return envelopes, nil
}

func (d *diskPrevalentStreamStore) Write(envelopes []*prevalentEventEnvelope) error {
	bytes, err := d.encode(envelopes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		fmt.Sprintf("%s/%s", d.folder, uuid.New().String()),
		bytes,
		0644,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *diskPrevalentStreamStore) encode(envelopes []*prevalentEventEnvelope) ([]byte, error) {
	encodedEnvelopes := []map[string]interface{}{}

	for _, envelope := range envelopes {
		event, err := json.Marshal(envelope.Event())
		if err != nil {
			return nil, err
		}

		encodedEnvelopes = append(encodedEnvelopes, map[string]interface{}{
			"event":          string(event),
			"eventType":      envelope.Event().Type(),
			"streamPosition": envelope.StreamPosition().Value(),
			"timestamp":      envelope.Timestamp(),
		})
	}

	return json.Marshal(encodedEnvelopes)
}

func (d *diskPrevalentStreamStore) decode(bytes []byte) ([]*prevalentEventEnvelope, error) {
	encodedEnvelopes := []map[string]*json.RawMessage{}
	err := json.Unmarshal(bytes, &encodedEnvelopes)
	if err != nil {
		return nil, err
	}

	envelopes := []*prevalentEventEnvelope{}

	for _, encodedEnvelope := range encodedEnvelopes {
		var timestamp time.Time
		err = json.Unmarshal(*encodedEnvelope["timestamp"], &timestamp)
		if err != nil {
			return nil, err
		}

		var streamPosition uint64
		err = json.Unmarshal(*encodedEnvelope["streamPosition"], &streamPosition)
		if err != nil {
			return nil, err
		}

		var eventType string
		err = json.Unmarshal(*encodedEnvelope["eventType"], &eventType)
		if err != nil {
			return nil, err
		}

		var event string
		err = json.Unmarshal(*encodedEnvelope["event"], &event)
		if err != nil {
			return nil, err
		}

		eventInstance, err := d.registry.NewEvent(eventType)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(event), &eventInstance)
		if err != nil {
			return nil, err
		}

		envelopes = append(envelopes, newPrevalentEventEnvelope(
			eventInstance,
			newPrevalentStreamPosition(streamPosition),
			timestamp,
		))
	}

	return envelopes, nil
}
