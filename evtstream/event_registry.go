package evtstream

import (
	"fmt"
	"reflect"
)

type EventRegistry interface {
	NewEvent(eventType string) (Event, error)
}

type inMemoryEventRegistry struct {
	eventMap map[string]Event
}

func NewInMemoryEventRegistry(events []Event) *inMemoryEventRegistry {
	eventMap := map[string]Event{}

	for _, event := range events {
		eventMap[event.Type()] = event
	}

	return &inMemoryEventRegistry{eventMap: eventMap}
}

func (i *inMemoryEventRegistry) NewEvent(eventType string) (Event, error) {
	if event, ok := i.eventMap[eventType]; ok {
		reflectType := reflect.ValueOf(event).Elem().Type()
		return reflect.New(reflectType).Interface().(Event), nil
	}

	return nil, fmt.Errorf("event %s not registred", eventType)
}
