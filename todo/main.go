package todo

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
)

//// COMMANDS

// CmdTaskCreate add task command
type CmdTaskCreate struct {
	ID          string
	Description string
}

//// EVENTS

// EvtTaskCreated task created event
type EvtTaskCreated struct {
	ID          string
	Description string
}

//// Projections

//// Domain Projection
type taskDomainProjection struct {
	id          string
	description string
}

func (m *taskDomainProjection) apply(eventEnvelope *eventstream.EventEnvelope) error {
	switch eventEnvelope.Type {
	case "TASK_CREATED":
		evt := &EvtTaskCreated{}
		err := json.Unmarshal(eventEnvelope.Event, evt)
		if err != nil {
			return nil
		}
		m.applyTaskCreated(evt)
	}

	return nil
}

func (m *taskDomainProjection) applyTaskCreated(evt *EvtTaskCreated) {
	m.id = evt.ID
	m.description = evt.Description
}

//// COMMAND HANDLER

// cmdHandler execute commands to model
type cmdHandler struct {
	eventStream eventstream.EventStream
}

// Handle handles any task Cmd*
func (c *cmdHandler) Handle(cmd interface{}) error {
	switch v := cmd.(type) {
	case *CmdTaskCreate:
		return c.handleCmdTaskCreate(v)
	default:
		return nil
	}
}

func (c *cmdHandler) handleCmdTaskCreate(cmd *CmdTaskCreate) error {
	taskProjection := &taskDomainProjection{}

	aggregateEvents, err := c.eventStream.ReadAggregate(cmd.ID)
	if err != nil {
		return nil
	}

	if len(aggregateEvents) > 0 {
		return nil
	}

	for _, evt := range aggregateEvents {
		err := taskProjection.apply(evt)
		if err != nil {
			return err
		}
	}

	event, err := json.Marshal(&EvtTaskCreated{
		ID:          cmd.ID,
		Description: cmd.Description,
	})

	if err != nil {
		return nil
	}

	events := []*eventstream.EventEnvelope{
		{
			Type:             "TASK_CREATED",
			AggregateID:      cmd.ID,
			AggregateType:    "TASK",
			AggregateVersion: 1,
			Event:            event,
		},
	}

	return c.eventStream.Write(events)
}

// NewCmdHandler func factory of cmdHandler
func NewCmdHandler(eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream: eventStream}
}
