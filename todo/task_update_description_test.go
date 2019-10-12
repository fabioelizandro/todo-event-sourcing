package todo_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/todo"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_task_update_description(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &todo.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.Nil(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := todo.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		createdEvent,
		&todo.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_update_description_to_same(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &todo.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.Nil(t, eventStream.Write([]eventstream.Event{createdEvent}))

	cmdHandler := todo.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))
	assert.Nil(t, cmdHandler.Handle(cmd))

	expectedEvents := []eventstream.Event{
		createdEvent,
		&todo.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}

func Test_task_update_description_not_found(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}

	cmdHandler := todo.NewCmdHandler(eventStream)
	assert.Nil(t, cmdHandler.Handle(cmd))

	assert.Equal(t, 0, len(eventStream.InMemoryReadAll()))
}
