package todo

import "fabioelizandro/todo-event-sourcing/eventstream"

// CmdTodoCreate Add todo command
type CmdTodoCreate struct {

}

// cmdHandler execute commands to model
type cmdHandler struct {
	eventStream eventstream.EventStream
}

// Handle handles any todo Cmd*
func (c *cmdHandler) Handle(cmd interface{}) error {
	return nil
}


type eventStore interface {

}

// NewCmdHandler func factory of cmdHandler
func NewCmdHandler (eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream:eventStream}
}