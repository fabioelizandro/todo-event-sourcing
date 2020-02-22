package task

type EvtTaskCreated struct {
	ID          string
	Description string
	Completed   bool
	CreatedAt   int64
}

func (e *EvtTaskCreated) Type() string {
	return "TASK_CREATED"
}

func (e *EvtTaskCreated) CorrelationID() string {
	return e.ID
}

type EvtTaskDescriptionUpdated struct {
	ID          string
	Description string
}

func (e *EvtTaskDescriptionUpdated) Type() string {
	return "TASK_DESCRIPTION_UPDATED"
}

func (e *EvtTaskDescriptionUpdated) CorrelationID() string {
	return e.ID
}

type EvtTaskCompleted struct {
	ID string
}

func (e *EvtTaskCompleted) Type() string {
	return "TASK_COMPLETED"
}

func (e *EvtTaskCompleted) CorrelationID() string {
	return e.ID
}
