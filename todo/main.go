package todo

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
)

//// COMMANDS

type CmdTaskCreate struct {
	ID          string
	Description string
}

//// EVENTS

type EvtTaskCreated struct {
	ID          string
	Description string
}

func (e *EvtTaskCreated) Type() string {
	return "TASK_CREATED"
}

func (e *EvtTaskCreated) AggregateID() string {
	return e.ID
}

func (e *EvtTaskCreated) AggregateType() string {
	return "TASK"
}

func (e *EvtTaskCreated) AggregateVersion() int64 {
	return 1
}

func (e *EvtTaskCreated) Payload() ([]byte, error) {
	return json.Marshal(e)
}

//// Projections

type taskDomainProjection struct {
	id          string
	description string
}

func (m *taskDomainProjection) apply(evt eventstream.Event) {
	switch v := evt.(type) {
	case *EvtTaskCreated:
		m.applyTaskCreated(v)
	}
}

func (m *taskDomainProjection) applyTaskCreated(evt *EvtTaskCreated) {
	m.id = evt.ID
	m.description = evt.Description
}

//// COMMAND HANDLER

type cmdHandler struct {
	eventStream eventstream.EventStream
}

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
		taskProjection.apply(evt)
	}

	events := []eventstream.Event{
		&EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}

	return c.eventStream.Write(events)
}

func NewCmdHandler(eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream: eventStream}
}
