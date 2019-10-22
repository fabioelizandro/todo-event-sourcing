package main

import (
	"encoding/json"
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type taskListRoute struct {
	projection taskprojection.TaskProjection
}

func (t *taskListRoute) Func(w http.ResponseWriter, r *http.Request) {
	tasks := t.projection.Tasks()
	b, err := json.Marshal(tasks)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}

type taskCreateRoute struct {
	commandHandler task.CmdHandler
}

func (t *taskCreateRoute) Func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload map[string]string
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	cmd := &task.CmdTaskCreate{
		ID:          uuid.New().String(),
		Description: payload["description"],
		CreatedAt:   time.Now().UnixNano(),
	}
	err = t.commandHandler.Handle(cmd)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"message":"OK"}`))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	stream := eventstream.NewInMemoryEventStream()
	commandHandler := task.NewCmdHandler(stream)
	projection := taskprojection.NewTaskProjection(stream)
	go projection.PollEventStream(100)

	r := mux.NewRouter()
	r.HandleFunc("/todos", (&taskListRoute{projection: projection}).Func).Methods("GET")
	r.HandleFunc("/todos", (&taskCreateRoute{commandHandler: commandHandler}).Func).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
