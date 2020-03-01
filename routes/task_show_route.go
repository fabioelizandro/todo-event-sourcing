package routes

import (
	"fabioelizandro/todo-event-sourcing/lib/http_essentials"
	"fabioelizandro/todo-event-sourcing/taskprojection"
)

type taskShowRoute struct {
	projection taskprojection.TaskProjection
}

func NewTaskShowRoute(projection taskprojection.TaskProjection) *taskShowRoute {
	return &taskShowRoute{projection: projection}
}

func (t *taskShowRoute) Methods() []string {
	return []string{"GET"}
}

func (t *taskShowRoute) Path() string {
	return "/todos/{id}"
}

func (t *taskShowRoute) Handle(r http_essentials.Request) (http_essentials.Response, error) {
	return http_essentials.NewJsonResponse(http_essentials.Headers{}, t.projection.Task(r.PathParams().Value("id", ""))), nil
}
