package task

import "encoding/json"

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
