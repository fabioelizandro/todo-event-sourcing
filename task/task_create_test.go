package task_test

import (
	"fabioelizandro/todo-event-sourcing/lib/evtstream"
	"fabioelizandro/todo-event-sourcing/task"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_it_create_tasks(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskCreate{
		ID:          uuid.New().String(),
		Description: "Do the dishes",
		CreatedAt:   time.Now().UnixNano(),
	}
	cmdHandler := task.NewCmdHandler(eventStream)

	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []evtstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
			CreatedAt:   cmd.CreatedAt,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.Tape())
}

func Test_it_requires_description_when_creating(t *testing.T) {
	cmd := &task.CmdTaskCreate{
		ID:          "123",
		Description: "",
		CreatedAt:   time.Now().UnixNano(),
	}
	cmdHandler := task.NewCmdHandler(evtstream.NewRecordingEventStream())

	rejection, err := cmdHandler.Handle(cmd)
	assert.NoError(t, err)

	assert.Equal(t, rejection, &task.CmdRejectionRequiredField{Name: "Description"})
}

func Test_it_does_not_create_tasks_with_same_id(t *testing.T) {
	eventStream := evtstream.NewRecordingEventStream()
	cmd := &task.CmdTaskCreate{ID: uuid.New().String(), Description: "Do the dishes"}
	cmdHandler := task.NewCmdHandler(eventStream)

	rejection, err := cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	rejection, err = cmdHandler.Handle(cmd)
	assert.Nil(t, rejection)
	assert.NoError(t, err)

	expectedEvents := []evtstream.Event{
		&task.EvtTaskCreated{
			ID:          cmd.ID,
			Description: cmd.Description,
		},
	}
	assert.Equal(t, expectedEvents, eventStream.Tape())
}
