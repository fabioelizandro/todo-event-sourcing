package todo

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
)

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
		return nil
	}
}

func (c *cmdHandler) handleCmdTaskCreate(cmd *CmdTaskCreate) error {
	domainModel, err := c.loadDomainModel(cmd.ID)
	if err != nil {
		return nil
	}

	events := domainModel.create(cmd.ID, cmd.Description)
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
	aggregateEvents, err := c.eventStream.ReadAggregate(ID)
	if err != nil {
		return nil, err
	}

	for _, evt := range aggregateEvents {
		domainModel.apply(evt)
	}

	return domainModel, nil
}
