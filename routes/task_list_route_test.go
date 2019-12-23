package routes_test

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/routes"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_returns_task_list_route_configuration(t *testing.T) {
	tasks := map[string]*taskprojection.Task{}
	route := routes.NewTaskListRoute(taskprojection.NewFakeTaskProjection(tasks))
	assert.Equal(t, []string{"GET"}, route.Methods())
	assert.Equal(t, "/todos", route.Path())
}

func Test_it_returns_tasks_projection_list(t *testing.T) {
	request := http_essentials.NewRequest(http_essentials.Headers{}, http_essentials.NewEmptyFakeRequestBody(), http_essentials.PathParams{})
	tasks := map[string]*taskprojection.Task{
		"1": {
			ID:          "1",
			Description: "Do something",
			Completed:   false,
			CreatedAt:   0,
		},
	}

	route := routes.NewTaskListRoute(taskprojection.NewFakeTaskProjection(tasks))
	response, err := route.Handle(request)
	assert.NoError(t, err)

	expectedBody, err := json.Marshal([]*taskprojection.Task{tasks["1"]})
	assert.NoError(t, err)

	actualBody, err := response.Body()
	assert.NoError(t, err)

	assert.Equal(t, "200", response.Headers().Value("status", ""))
	assert.Equal(t, "application/json", response.Headers().Value("content-type", ""))
	assert.Equal(t, string(expectedBody), string(actualBody))
}
