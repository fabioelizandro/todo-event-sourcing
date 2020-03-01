package routes_test

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/lib/http_essentials"
	"fabioelizandro/todo-event-sourcing/routes"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_returns_task_show_route_configuration(t *testing.T) {
	tasks := map[string]*taskprojection.Task{}
	route := routes.NewTaskShowRoute(taskprojection.NewFakeTaskProjection(tasks))
	assert.Equal(t, []string{"GET"}, route.Methods())
	assert.Equal(t, "/todos/{id}", route.Path())
}

func Test_it_returns_tasks_projection_show(t *testing.T) {
	request := http_essentials.NewRequest(http_essentials.Headers{}, http_essentials.NewEmptyFakeRequestBody(), http_essentials.PathParams{
		"id": "1",
	})
	tasks := map[string]*taskprojection.Task{
		"1": {
			ID:          "1",
			Description: "Do something",
			Completed:   false,
			CreatedAt:   0,
		},
	}

	route := routes.NewTaskShowRoute(taskprojection.NewFakeTaskProjection(tasks))
	response, err := route.Handle(request)
	assert.NoError(t, err)

	expectedBody, err := json.Marshal(tasks["1"])
	assert.NoError(t, err)

	actualBody, err := response.Body()
	assert.NoError(t, err)

	assert.Equal(t, "200", response.Headers().Value("status", ""))
	assert.Equal(t, "application/json", response.Headers().Value("content-type", ""))
	assert.Equal(t, string(expectedBody), string(actualBody))
}
