package taskprojections_test

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojections"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_lists_all_created_tasks(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.Nil(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", CreatedAt: 0},
		&task.EvtTaskCreated{ID: "456", Description: "Clean house", CreatedAt: 1},
	}))

	listProjection := taskprojections.NewTaskListProjection(eventStream)
	assert.Nil(t, listProjection.CatchupEventStream())

	expectedTasks := []*taskprojections.Task{
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

	assert.Equal(t, expectedTasks, listProjection.Tasks())
}

func Test_it_updates_task_descriptions(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.Nil(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes"},
		&task.EvtTaskDescriptionUpdated{ID: "123", Description: "Clean house"},
	}))

	listProjection := taskprojections.NewTaskListProjection(eventStream)
	assert.Nil(t, listProjection.CatchupEventStream())

	expectedTasks := []*taskprojections.Task{
		{
			ID:          "123",
			Description: "Clean house",
			Completed:   false,
		},
	}

	assert.Equal(t, expectedTasks, listProjection.Tasks())
}

func Test_it_marks_task_as_completed(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.Nil(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", Completed: false},
		&task.EvtTaskCompleted{ID: "123"},
	}))

	listProjection := taskprojections.NewTaskListProjection(eventStream)
	assert.Nil(t, listProjection.CatchupEventStream())

	expectedTasks := []*taskprojections.Task{
		{
			ID:          "123",
			Description: "Do the dishes",
			Completed:   true,
		},
	}

	assert.Equal(t, expectedTasks, listProjection.Tasks())
}

func Test_it_shows_task_by_id(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	assert.Nil(t, eventStream.Write([]eventstream.Event{
		&task.EvtTaskCreated{ID: "123", Description: "Do the dishes", Completed: false},
	}))

	listProjection := taskprojections.NewTaskListProjection(eventStream)
	assert.Nil(t, listProjection.CatchupEventStream())

	expectedTask := &taskprojections.Task{
		ID:          "123",
		Description: "Do the dishes",
		Completed:   false,
	}

	assert.Equal(t, expectedTask, listProjection.Task("123"))
}

func Test_it_returns_nil_when_task_not_found(t *testing.T) {
	eventStream := eventstream.NewInMemoryEventStream()
	listProjection := taskprojections.NewTaskListProjection(eventStream)
	assert.Nil(t, listProjection.Task("123"))
}
