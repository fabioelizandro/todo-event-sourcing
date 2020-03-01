package task

import "fabioelizandro/todo-event-sourcing/lib/evtstream"

type taskDomainModel struct {
	id          string
	description string
	completed   bool
}

func (m *taskDomainModel) apply(evtEnvelope evtstream.EventEnvelope) {
	evt := evtEnvelope.Event()

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

func (m *taskDomainModel) updateDescription(newDescription string) ([]evtstream.Event, CmdRejection) {
	events := make([]evtstream.Event, 0)

	if m.id == "" {
		return events, nil
	}

	if len(newDescription) == 0 {
		return nil, &CmdRejectionRequiredField{Name: "NewDescription"}
	}

	if m.description != newDescription {
		events = append(events, &EvtTaskDescriptionUpdated{
			ID:          m.id,
			Description: newDescription,
		})
	}

	return events, nil
}

func (m *taskDomainModel) complete() []evtstream.Event {
	events := make([]evtstream.Event, 0)

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

func (m *taskDomainModel) create(ID string, description string, createdAt int64) ([]evtstream.Event, CmdRejection) {
	events := make([]evtstream.Event, 0)

	if m.id != "" {
		return events, nil
	}

	if len(description) == 0 {
		return nil, &CmdRejectionRequiredField{Name: "Description"}
	}

	events = []evtstream.Event{
		&EvtTaskCreated{
			ID:          ID,
			Description: description,
			CreatedAt:   createdAt,
		},
	}

	return events, nil
}
