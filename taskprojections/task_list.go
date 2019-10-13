package taskprojections

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"sort"
	"time"
)

type Task struct {
	ID          string
	Description string
	Completed   bool
	CreatedAt   int64
}

type taskListProjection struct {
	es      eventstream.EventStream
	eventID uint64
	tasks   map[string]*Task
}

func (t *taskListProjection) apply(evt eventstream.Event) {
	switch v := evt.(type) {
	case *task.EvtTaskCreated:
		t.applyTaskCreated(v)
	case *task.EvtTaskDescriptionUpdated:
		t.applyTaskDescriptionUpdated(v)
	case *task.EvtTaskCompleted:
		t.applyTaskCompleted(v)
	}
}

func (t *taskListProjection) applyTaskCreated(evt *task.EvtTaskCreated) {
	t.tasks[evt.ID] = &Task{
		ID:          evt.ID,
		Description: evt.Description,
		Completed:   evt.Completed,
	}
}

func (t *taskListProjection) applyTaskDescriptionUpdated(evt *task.EvtTaskDescriptionUpdated) {
	t.tasks[evt.ID].Description = evt.Description
}

func (t *taskListProjection) applyTaskCompleted(evt *task.EvtTaskCompleted) {
	t.tasks[evt.ID].Completed = true
}

func (t *taskListProjection) CatchupEventStream() error {
	for {
		evt, err := t.es.Read(t.eventID)
		if err != nil {
			return err
		}

		if evt == nil {
			return nil
		}

		t.apply(evt)
		t.eventID++
	}
}

func (t *taskListProjection) PollEventStream(intervalMilliseconds int) error {
	for {
		evt, err := t.es.Read(t.eventID)
		if err != nil {
			return err
		}

		t.apply(evt)
		t.eventID++

		time.Sleep(time.Duration(intervalMilliseconds) * time.Millisecond)
	}
}

func (t *taskListProjection) Tasks() []*Task {
	tasks := make([]*Task, 0)

	for _, v := range t.tasks {
		tasks = append(tasks, v)
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt < tasks[j].CreatedAt
	})

	return tasks
}

func (t *taskListProjection) Task(ID string) *Task {
	return t.tasks[ID]
}

func NewTaskListProjection(es eventstream.EventStream) *taskListProjection {
	projection := &taskListProjection{es: es}
	projection.tasks = make(map[string]*Task, 0)
	return projection
}
