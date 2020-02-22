package taskprojection

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"sort"
)

type Task struct {
	ID          string
	Description string
	Completed   bool
	CreatedAt   int64
}

type TaskProjection interface {
	Tasks() []*Task
	Task(ID string) *Task
	CatchupEventStream() error
}

type fakeTaskProjection struct {
	tasks map[string]*Task
}

func (f *fakeTaskProjection) Tasks() []*Task {
	tasks := make([]*Task, 0)
	for _, task := range f.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (f *fakeTaskProjection) Task(ID string) *Task {
	return f.tasks[ID]
}

func (f *fakeTaskProjection) CatchupEventStream() error {
	return nil
}

type taskProjection struct {
	es             eventstream.EventStream
	streamPosition eventstream.StreamPosition
	tasks          map[string]*Task
}

func (t *taskProjection) apply(evt eventstream.Event) {
	switch v := evt.(type) {
	case *task.EvtTaskCreated:
		t.applyTaskCreated(v)
	case *task.EvtTaskDescriptionUpdated:
		t.applyTaskDescriptionUpdated(v)
	case *task.EvtTaskCompleted:
		t.applyTaskCompleted(v)
	}
}

func (t *taskProjection) applyTaskCreated(evt *task.EvtTaskCreated) {
	t.tasks[evt.ID] = &Task{
		ID:          evt.ID,
		Description: evt.Description,
		Completed:   evt.Completed,
	}
}

func (t *taskProjection) applyTaskDescriptionUpdated(evt *task.EvtTaskDescriptionUpdated) {
	t.tasks[evt.ID].Description = evt.Description
}

func (t *taskProjection) applyTaskCompleted(evt *task.EvtTaskCompleted) {
	t.tasks[evt.ID].Completed = true
}

func (t *taskProjection) CatchupEventStream() error {
	for {
		result, err := t.es.Read(t.streamPosition)
		if err != nil {
			return err
		}

		if result == nil {
			return nil
		}

		t.apply(result.Event())
		t.streamPosition = result.NextStreamPosition()
	}
}

func (t *taskProjection) Tasks() []*Task {
	tasks := make([]*Task, 0)

	for _, v := range t.tasks {
		tasks = append(tasks, v)
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt < tasks[j].CreatedAt
	})

	return tasks
}

func (t *taskProjection) Task(ID string) *Task {
	return t.tasks[ID]
}

func NewTaskProjection(es eventstream.EventStream) *taskProjection {
	return &taskProjection{
		es:             es,
		tasks:          make(map[string]*Task, 0),
		streamPosition: es.FirstPosition(),
	}
}

func NewFakeTaskProjection(tasks map[string]*Task) *fakeTaskProjection {
	return &fakeTaskProjection{tasks: tasks}
}
