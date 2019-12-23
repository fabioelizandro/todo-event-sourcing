package main

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/http_routes"
	"fabioelizandro/todo-event-sourcing/logger"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	zlog := logger.NewZLog()
	stream := eventstream.NewInMemoryEventStream()
	commandHandler := task.NewCmdHandler(stream)
	projection := taskprojection.NewTaskProjection(stream)
	routeAdapter := http_routes.NewStdHttpRouteAdapter(zlog)
	routes := []http_routes.Route{
		http_routes.NewTaskListRoute(projection),
		http_routes.NewTaskShowRoute(projection),
		http_routes.NewTaskCreateRoute(commandHandler),
	}

	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.Path(), routeAdapter.Transform(route)).Methods(route.Methods()...)
	}

	go PollEventStream(projection)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func PollEventStream(projection taskprojection.TaskProjection) {
	for {
		err := projection.CatchupEventStream()
		if err != nil {
			fmt.Printf("TaskProjection error: %e", err)
		}

		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
