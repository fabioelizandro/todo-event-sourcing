package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_task_created(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := task.NewCmdHandler(eventStream)

	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_create_duplicated(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := task.NewCmdHandler(eventStream)

	assert.Nil(t, cmdHandler.Handle(cmd))
	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
