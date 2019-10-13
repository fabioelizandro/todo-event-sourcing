package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_task_complete(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.Nil(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		createdEvent,
		&task.EvtTaskCompleted{ID: cmd.ID},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_complete_already_completed(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.Nil(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))
	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		createdEvent,
		&task.EvtTaskCompleted{ID: cmd.ID},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_complete_not_found(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}

	cmdHandler := task.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))

	assert.Equal(t, 0, len(eventStream.InMemoryReadAll()))
}
