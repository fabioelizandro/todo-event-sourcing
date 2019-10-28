package task_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_returns_an_error_when_cmd_not_found(t *testing.T) {
	type CmdNotFound struct{}
	cmdHandler := task.NewCmdHandler(eventstream.NewInMemoryEventStream())

	cmd := &CmdNotFound{}
	_, err := cmdHandler.Handle(cmd)

	assert.Equal(t, fmt.Errorf("command not found %v", cmd), err)
}
