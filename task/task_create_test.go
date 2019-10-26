package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_create_tasks(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskCreate{
		ID:          uuid.New().String(),
		Description: "Do the dishes",
		CreatedAt:   time.Now().UnixNano(),
	}
	cmdHandler := task.NewCmdHandler(eventStream)

	assert.NoError(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
			CreatedAt:   cmd.CreatedAt,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_it_does_not_create_tasks_with_same_id(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := task.NewCmdHandler(eventStream)

	assert.NoError(t, cmdHandler.Handle(cmd))
	assert.NoError(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
