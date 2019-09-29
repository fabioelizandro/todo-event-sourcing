package todo_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/todo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddTodo(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	cmd := &todo.CmdTodoCreate{}
	cmdHandler := todo.NewCmdHandler(eventStream)
	err := cmdHandler.Handle(cmd)

	expectedEvents := []eventstream.Event{
		{
			ID:               "1",
			Type:             "TODO_CREATED",
			AggregateID:      "1",
			AggregateType:    "TODO",
			AggregateVersion: 1,
			Payload:          "?",
		},
	}
	assert.Nil(t, err)
	assert.Equal(t, expectedEvents, eventStream.InMemoryReadAll())
}
