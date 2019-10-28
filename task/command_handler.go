package task

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fmt"
)

type CmdHandler interface {
	Handle(cmd Cmd) (CmdRejection, error)
}

type fakeCmdHandler struct {
	executedCmds []interface{}
}

func NewFakeCmdHandler() *fakeCmdHandler {
	return &fakeCmdHandler{}
}

func (f *fakeCmdHandler) Handle(cmd Cmd) (CmdRejection, error) {
	f.executedCmds = append(f.executedCmds, cmd)
	return nil, nil
}

func (f *fakeCmdHandler) ExecutedCmds() []interface{} {
	return f.executedCmds
}

type cmdHandler struct {
	eventStream eventstream.EventStream
}

func NewCmdHandler(eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream: eventStream}
}

func (c *cmdHandler) Handle(cmd Cmd) (CmdRejection, error) {
	switch v := cmd.(type) {
	case *CmdTaskCreate:
		return c.handleCmdTaskCreate(v)
	case *CmdTaskUpdateDescription:
		return c.handleCmdTaskUpdateDescription(v)
	case *CmdTaskComplete:
		return c.handleCmdTaskComplete(v)
	default:
		return nil, fmt.Errorf("command not found %v", v)
	}
}

func (c *cmdHandler) handleCmdTaskCreate(cmd *CmdTaskCreate) (CmdRejection, error) {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return nil, err
	}

	events, rejection := domainModel.create(cmd.ID, cmd.Description, cmd.CreatedAt)
	if rejection != nil {
		return rejection, nil
	}

	return nil, c.eventStream.Write(events)
}

func (c *cmdHandler) handleCmdTaskUpdateDescription(cmd *CmdTaskUpdateDescription) (CmdRejection, error) {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return nil, err
	}

	events := domainModel.updateDescription(cmd.NewDescription)
	return nil, c.eventStream.Write(events)
}

func (c *cmdHandler) handleCmdTaskComplete(cmd *CmdTaskComplete) (CmdRejection, error) {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return nil, err
	}

	events := domainModel.complete()
	return nil, c.eventStream.Write(events)
}

func (c *cmdHandler) loadDomainModel(ID string) (*taskDomainModel, error) {
	domainModel := &taskDomainModel{}
	aggregateEvents, err := c.eventStream.ReadAllByCorrelationID(ID)
	if err != nil {
		return nil, err
	}

	for _, evt := range aggregateEvents {
		domainModel.apply(evt)
	}

	return domainModel, nil
}
