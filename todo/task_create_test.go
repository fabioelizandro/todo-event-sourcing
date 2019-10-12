package todo_test

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/todo"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_task_created(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := todo.NewCmdHandler(eventStream)

	assert.Nil(t, cmdHandler.Handle(cmd))

	event, err := json.Marshal(&todo.EvtTaskCreated{
		ID:          cmd.ID,
		Description: cmd.Description,
	})
	assert.Nil(t, err)

	expectedEvents := []*eventstream.EventEnvelope{
		{
			Type:             "TASK_CREATED",
			AggregateID:      cmd.ID,
			AggregateType:    "TASK",
			AggregateVersion: 1,
			Event:            event,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_create_duplicated(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := todo.NewCmdHandler(eventStream)

	assert.Nil(t, cmdHandler.Handle(cmd))
	assert.Nil(t, cmdHandler.Handle(cmd))

	event, err := json.Marshal(&todo.EvtTaskCreated{
		ID:          cmd.ID,
		Description: cmd.Description,
	})
	assert.Nil(t, err)

	expectedEvents := []*eventstream.EventEnvelope{
		{
			Type:             "TASK_CREATED",
			AggregateID:      cmd.ID,
			AggregateType:    "TASK",
			AggregateVersion: 1,
			Event:            event,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
