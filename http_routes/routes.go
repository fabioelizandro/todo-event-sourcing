package http_routes

import (
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"time"

	"github.com/google/uuid"
)

type Route interface {
	Methods() []string
	Path() string
	Handle(Request) (Response, error)
}

type taskListRoute struct {
	projection taskprojection.TaskProjection
}

func (t *taskListRoute) Methods() []string {
	return []string{"GET"}
}

func (t *taskListRoute) Path() string {
	return "/todos"
}

func (t *taskListRoute) Handle(Request) (Response, error) {
	return newJsonResponse(Headers{}, t.projection.Tasks()), nil
}

type taskCreateRoute struct {
	commandHandler task.CmdHandler
}

func (t *taskCreateRoute) Methods() []string {
	return []string{"POST"}
}

func (t *taskCreateRoute) Path() string {
	return "/todos"
}

func (t *taskCreateRoute) Handle(r Request) (Response, error) {
	cmd := &task.CmdTaskCreate{
		ID:          uuid.New().String(),
		Description: r.Body().FieldStr("description", ""),
		CreatedAt:   time.Now().UnixNano(),
	}

	err := t.commandHandler.Handle(cmd)
	if err != nil {
		return nil, err
	}

	return newJsonOkResponse(Headers{}), nil
}

func NewTaskListRoute(projection taskprojection.TaskProjection) *taskListRoute {
	return &taskListRoute{projection: projection}
}

func NewTaskCreateRoute(commandHandler task.CmdHandler) *taskCreateRoute {
	return &taskCreateRoute{commandHandler: commandHandler}
}
