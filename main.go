package main

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/logger"
	"fabioelizandro/todo-event-sourcing/routes"
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
	routeAdapter := http_essentials.NewStdHttpRouteAdapter(zlog)
	httpRoutes := []http_essentials.Route{
		routes.NewTaskListRoute(projection),
		routes.NewTaskShowRoute(projection),
		routes.NewTaskCreateRoute(commandHandler),
	}

	r := mux.NewRouter()
	for _, route := range httpRoutes {
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
