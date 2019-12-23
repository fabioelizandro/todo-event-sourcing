package routes_test

import (
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/routes"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_returns_task_create_route_configuration(t *testing.T) {
	route := routes.NewTaskCreateRoute(task.NewFakeCmdHandler())
	assert.Equal(t, []string{"POST"}, route.Methods())
	assert.Equal(t, "/todos", route.Path())
}

func Test_it_execute_create_cmd_handler(t *testing.T) {
	cmdHandler := task.NewFakeCmdHandler()
	requestBody := http_essentials.NewFakeRequestBody(http_essentials.RequestBodyFields{
		"description": "Do the dishes",
	})
	request := http_essentials.NewRequest(http_essentials.Headers{}, requestBody, http_essentials.PathParams{})

	route := routes.NewTaskCreateRoute(cmdHandler)
	response, err := route.Handle(request)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(cmdHandler.ExecutedCmds()))
	assert.IsType(t, &task.CmdTaskCreate{}, cmdHandler.ExecutedCmds()[0])
	assert.Equal(t, "202", response.Headers().Value("status", ""))
	assert.Equal(t, "application/json", response.Headers().Value("content-type", ""))
}
