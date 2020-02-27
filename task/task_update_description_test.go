package task_test

import (
	"fabioelizandro/todo-event-sourcing/evtstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_updates_task_description(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.NoError(t, eventStream.Write([]evtstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []evtstream.Event{
		createdEvent,
		&task.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.Tape())
}

func Test_it_requires_description_when_updating(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: ""}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.NoError(t, eventStream.Write([]evtstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.NoError(t, err)

	assert.Equal(t, rejection, &task.CmdRejectionRequiredField{Name: "NewDescription"})
}

func Test_it_ignores_cmd_when_new_description_is_equal_to_current(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}
	createdEvent := &task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"}
	assert.NoError(t, eventStream.Write([]evtstream.Event{createdEvent}))

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	rejection, err = cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []evtstream.Event{
		createdEvent,
		&task.EvtTaskDescriptionUpdated{
			ID:          cmd.ID,
			Description: cmd.NewDescription,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.Tape())
}

func Test_it_ignores_cmd_when_task_is_not_found(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskUpdateDescription{ID: uuid.New().String(), NewDescription: "Clean kitchen"}

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(eventStream.Tape()))
}
