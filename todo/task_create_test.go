package todo_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/todo"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreate(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := todo.NewCmdHandler(eventStream)

	err := cmdHandler.Handle(cmd)

	expectedEvents := []*eventstream.EventEnvelope{
		{
			Type:             "TASK_CREATED",
			AggregateID:      cmd.ID,
			AggregateType:    "TASK",
			AggregateVersion: 1,
			Event: &todo.EvtTaskCreated{
				ID:          cmd.ID,
				Description: cmd.Description,
			},
		},
	}
	assert.Nil(t, err)
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
