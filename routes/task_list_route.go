package routes

import (
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/taskprojection"
)

type taskListRoute struct {
	projection taskprojection.TaskProjection
}

func NewTaskListRoute(projection taskprojection.TaskProjection) *taskListRoute {
	return &taskListRoute{projection: projection}
}

func (t *taskListRoute) Methods() []string {
	return []string{"GET"}
}

func (t *taskListRoute) Path() string {
	return "/todos"
}

func (t *taskListRoute) Handle(http_essentials.Request) (http_essentials.Response, error) {
	return http_essentials.NewJsonResponse(http_essentials.Headers{}, t.projection.Tasks()), nil
}
