package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_marks_task_as_completed(t *testing.T) {
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

func Test_it_ignores_complete_cmd_when_is_complete_already(t *testing.T) {
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

func Test_it_ignores_cmd_complete_for_not_found_tasks(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}

	cmdHandler := task.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))

	assert.Equal(t, 0, len(eventStream.InMemoryReadAll()))
}
