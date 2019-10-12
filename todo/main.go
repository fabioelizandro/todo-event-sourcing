package todo

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
)

// CmdTaskCreate add task command
type CmdTaskCreate struct {
	ID          string
	Description string
}

// EvtTaskCreated task created event
type EvtTaskCreated struct {
	ID          string
	Description string
}

func (e *EvtTaskCreated) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// cmdHandler execute commands to model
type cmdHandler struct {
	eventStream eventstream.EventStream
}

// Handle handles any task Cmd*
func (c *cmdHandler) Handle(cmd interface{}) error {
	switch v := cmd.(type) {
	case *CmdTaskCreate:
		return c.handleCmdTaskCreate(v)
	}

	return nil
}

func (c *cmdHandler) handleCmdTaskCreate(cmd *CmdTaskCreate) error {
	taskCreatedEvent := &EvtTaskCreated{
		ID:          cmd.ID,
		Description: cmd.Description,
	}

	events := []*eventstream.EventEnvelope{
		{
			Type:             "TASK_CREATED",
			AggregateID:      cmd.ID,
			AggregateType:    "TASK",
			AggregateVersion: 1,
			Event:            taskCreatedEvent,
		},
	}

	return c.eventStream.Write(events)
}

// NewCmdHandler func factory of cmdHandler
func NewCmdHandler(eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream: eventStream}
}
