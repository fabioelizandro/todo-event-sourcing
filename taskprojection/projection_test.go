package taskprojection_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_lists_all_created_tasks(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.NoError(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", CreatedAt: 0},
		&task.EvtTaskCreated{ID: "456", Description: "Clean house", CreatedAt: 1},
	}))

	projection := taskprojection.NewTaskProjection(eventStream)
	assert.NoError(t, projection.CatchupEventStream())

	expectedTasks := []*taskprojection.Task{
		{
			ID:          "123",
			Description: "Do the dishes",
			Completed:   false,
		},
		{
			ID:          "456",
			Description: "Clean house",
			Completed:   false,
		},
	}

	assert.Equal(t, expectedTasks, projection.Tasks())
}

func Test_it_updates_task_descriptions(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.NoError(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes"},
		&task.EvtTaskDescriptionUpdated{ID: "123", Description: "Clean house"},
	}))

	projection := taskprojection.NewTaskProjection(eventStream)
	assert.NoError(t, projection.CatchupEventStream())

	expectedTasks := []*taskprojection.Task{
		{
			ID:          "123",
			Description: "Clean house",
			Completed:   false,
		},
	}

	assert.Equal(t, expectedTasks, projection.Tasks())
}

func Test_it_marks_task_as_completed(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.NoError(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", Completed: false},
		&task.EvtTaskCompleted{ID: "123"},
	}))

	projection := taskprojection.NewTaskProjection(eventStream)
	assert.NoError(t, projection.CatchupEventStream())

	expectedTasks := []*taskprojection.Task{
		{
			ID:          "123",
			Description: "Do the dishes",
			Completed:   true,
		},
	}

	assert.Equal(t, expectedTasks, projection.Tasks())
}

func Test_it_shows_task_by_id(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.NoError(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", Completed: false},
	}))

	projection := taskprojection.NewTaskProjection(eventStream)
	assert.NoError(t, projection.CatchupEventStream())

	expectedTask := &taskprojection.Task{
		ID:          "123",
		Description: "Do the dishes",
		Completed:   false,
	}

	assert.Equal(t, expectedTask, projection.Task("123"))
}

func Test_it_returns_nil_when_task_not_found(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	projection := taskprojection.NewTaskProjection(eventStream)
	assert.Nil(t, projection.Task("123"))
}
