package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_updates_task_description(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.NoError(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []eventstream.Event{
		createdEvent,
		&task.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_it_ignores_cmd_when_new_description_is_equal_to_current(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.NoError(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	rejection, err = cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []eventstream.Event{
		createdEvent,
		&task.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_it_ignores_cmd_when_task_is_not_found(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(eventStream.InMemoryReadAll()))
}
