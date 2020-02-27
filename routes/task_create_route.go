package routes

import (
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/task"
	"time"

	"github.com/google/uuid"
)

type taskCreateRoute struct {
	commandHandler task.CmdHandler
}

func NewTaskCreateRoute(commandHandler task.CmdHandler) *taskCreateRoute {
	return &taskCreateRoute{commandHandler: commandHandler}
}

func (t *taskCreateRoute) Methods() []string {
	return []string{"POST"}
}

func (t *taskCreateRoute) Path() string {
	return "/todos"
}

func (t *taskCreateRoute) Handle(r http_essentials.Request) (http_essentials.Response, error) {
	cmd := &task.CmdTaskCreate{
		ID:          uuid.New().String(),
		Description: r.Body().FieldStr("description", ""),
		CreatedAt:   time.Now().UnixNano(),
	}

	rejection, err := t.commandHandler.Handle(cmd)
	if err != nil {
		return nil, err
	}

	if rejection != nil {
		return http_essentials.NewJsonResponse(http_essentials.Headers{
			"status": "400",
		}, rejection), nil
	}

	return http_essentials.NewJsonOkResponse(http_essentials.Headers{}), nil
}
