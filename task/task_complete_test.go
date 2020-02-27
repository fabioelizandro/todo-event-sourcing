package task_test

import (
	"fabioelizandro/todo-event-sourcing/evtstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_marks_task_as_completed(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}
	assert.NoError(t, eventStream.Write([]evtstream.Event{
		&task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"},
	}))
	eventStream.EraseTape()

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	assert.Equal(t, []evtstream.Event{
		&task.EvtTaskCompleted{ID: cmd.ID},
	}, eventStream.Tape())
}

func Test_it_ignores_complete_cmd_when_is_complete_already(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}
	assert.NoError(t, eventStream.Write([]evtstream.Event{
		&task.EvtTaskCreated{ID: cmd.ID, Description: "Do the dishes"},
	}))
	eventStream.EraseTape()

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	rejection, err = cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	assert.Equal(t, []evtstream.Event{
		&task.EvtTaskCompleted{ID: cmd.ID},
	}, eventStream.Tape())
}

func Test_it_ignores_cmd_complete_for_not_found_tasks(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskComplete{ID: uuid.New().String()}

	cmdHandler := task.NewCmdHandler(eventStream)
	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(eventStream.Tape()))
}
