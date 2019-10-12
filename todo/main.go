package todo

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
)

//// COMMANDS

type CmdTaskCreate struct {
	ID          string
	Description string
}

type CmdTaskUpdateDescription struct {
	ID             string
	NewDescription string
}

type CmdTaskComplete struct {
	ID string
}

//// EVENTS

type EvtTaskCreated struct {
	ID          string
	Description string
	Completed   bool
}

func (e *EvtTaskCreated) Type() string {
	return "TASK_CREATED"
}

func (e *EvtTaskCreated) AggregateID() string {
	return e.ID
}

func (e *EvtTaskCreated) AggregateType() string {
	return "TASK"
}

func (e *EvtTaskCreated) Payload() ([]byte, error) {
	return json.Marshal(e)
}

type EvtTaskDescriptionUpdated struct {
	ID          string
	Description string
}

func (e *EvtTaskDescriptionUpdated) Type() string {
	return "TASK_DESCRIPTION_UPDATED"
}

func (e *EvtTaskDescriptionUpdated) AggregateID() string {
	return e.ID
}

func (e *EvtTaskDescriptionUpdated) AggregateType() string {
	return "TASK"
}

func (e *EvtTaskDescriptionUpdated) Payload() ([]byte, error) {
	return json.Marshal(e)
}

type EvtTaskCompleted struct {
	ID          string
	Description string
}

func (e *EvtTaskCompleted) Type() string {
	return "TASK_COMPLETED"
}

func (e *EvtTaskCompleted) AggregateID() string {
	return e.ID
}

func (e *EvtTaskCompleted) AggregateType() string {
	return "TASK"
}

func (e *EvtTaskCompleted) Payload() ([]byte, error) {
	return json.Marshal(e)
}

//// Domain model projection

type taskDomainModel struct {
	id          string
	description string
	completed   bool
}

func (m *taskDomainModel) apply(evt eventstream.Event) {
	switch v := evt.(type) {
	case *EvtTaskCreated:
		m.applyTaskCreated(v)
	case *EvtTaskDescriptionUpdated:
		m.applyTaskDescriptionUpdated(v)
	case *EvtTaskCompleted:
		m.applyTaskCompleted(v)
	}
}

func (m *taskDomainModel) applyTaskCreated(evt *EvtTaskCreated) {
	m.id = evt.ID
	m.description = evt.Description
	m.completed = evt.Completed
}

func (m *taskDomainModel) applyTaskDescriptionUpdated(evt *EvtTaskDescriptionUpdated) {
	m.description = evt.Description
}

func (m *taskDomainModel) applyTaskCompleted(evt *EvtTaskCompleted) {
	m.completed = true
}

func (m *taskDomainModel) updateDescription(newDescription string) []eventstream.Event {
	events := make([]eventstream.Event, 0)

	if m.id == "" {
		return events
	}

	if m.description != newDescription {
		events = append(events, &EvtTaskDescriptionUpdated{
			ID:          m.id,
			Description: newDescription,
		})
	}

	return events
}

func (m *taskDomainModel) complete() []eventstream.Event {
	events := make([]eventstream.Event, 0)

	if m.id == "" {
		return events
	}

	if !m.completed {
		events = append(events, &EvtTaskCompleted{
			ID: m.id,
		})
	}

	return events
}

//// COMMAND HANDLER

type cmdHandler struct {
	eventStream eventstream.EventStream
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
	domainModel := &taskDomainModel{}

	aggregateEvents, err := c.eventStream.ReadAggregate(cmd.ID)
	if err != nil {
		return nil
	}

	for _, evt := range aggregateEvents {
		domainModel.apply(evt)
	}

	if len(aggregateEvents) > 0 {
		return nil
	}

	events := []eventstream.Event{
		&EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}

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

func NewCmdHandler(eventStream eventstream.EventStream) *cmdHandler {
	return &cmdHandler{eventStream: eventStream}
}
