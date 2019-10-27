package task

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fmt"
)

type CmdHandler interface {
	Handle(cmd interface{}) error
}

type fakeCmdHandler struct {
	executedCmds []interface{}
}

func NewFakeCmdHandler() *fakeCmdHandler {
	return &fakeCmdHandler{}
}

func (f *fakeCmdHandler) Handle(cmd interface{}) error {
	f.executedCmds = append(f.executedCmds, cmd)
	return nil
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

func (c *cmdHandler) Handle(cmd interface{}) error {
	switch v := cmd.(type) {
	case *CmdTaskCreate:
		return c.handleCmdTaskCreate(v)
	case *CmdTaskUpdateDescription:
		return c.handleCmdTaskUpdateDescription(v)
	case *CmdTaskComplete:
		return c.handleCmdTaskComplete(v)
	default:
		return fmt.Errorf("command not found %v", v)
	}
}

func (c *cmdHandler) handleCmdTaskCreate(cmd *CmdTaskCreate) error {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return err
	}

	events := domainModel.create(cmd.ID, cmd.Description, cmd.CreatedAt)
	return c.eventStream.Write(events)
}

func (c *cmdHandler) handleCmdTaskUpdateDescription(cmd *CmdTaskUpdateDescription) error {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return err
	}

	events := domainModel.updateDescription(cmd.NewDescription)
	return c.eventStream.Write(events)
}

func (c *cmdHandler) handleCmdTaskComplete(cmd *CmdTaskComplete) error {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return err
	}

	events := domainModel.complete()
	return c.eventStream.Write(events)
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
