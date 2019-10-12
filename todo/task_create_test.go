package todo_test

import (
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

	expectedEvents := []eventstream.Event{
		&todo.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
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

	expectedEvents := []eventstream.Event{
		&todo.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
