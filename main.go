package main

import (
	"fabioelizandro/todo-event-sourcing/evtstream"
	"fabioelizandro/todo-event-sourcing/http_essentials"
	"fabioelizandro/todo-event-sourcing/logger"
	"fabioelizandro/todo-event-sourcing/routes"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	zlog := logger.NewZLog()
	stream := loadStream()
	commandHandler := task.NewCmdHandler(stream)
	projection := taskprojection.NewTaskProjection(stream)
	routeAdapter := newRouteAdapter(zlog)
	httpRoutes := []http_essentials.Route{
		routes.NewTaskListRoute(projection),
		routes.NewTaskShowRoute(projection),
		routes.NewTaskCreateRoute(commandHandler),
	}

	r := mux.NewRouter()
	for _, route := range httpRoutes {
		r.HandleFunc(route.Path(), routeAdapter.Transform(route)).Methods(route.Methods()...)
	}

	go pollEventStream(projection, zlog)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadStream() evtstream.EventStream {
	store := evtstream.NewDiskPrevalentEventStore(
		"/tmp/todo-event-sourcing-stream",
		evtstream.NewInMemoryEventRegistry([]evtstream.Event{
			&task.EvtTaskCompleted{},
			&task.EvtTaskCreated{},
			&task.EvtTaskDescriptionUpdated{},
		}),
	)

	envelopes, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}

	return evtstream.NewPrevalentEventStream(store, envelopes, evtstream.NewUTCCLock())
}

func newRouteAdapter(log logger.Log) http_essentials.StdHttpRouteAdapter {
	return http_essentials.NewStdHttpRouteAdapter(func(msg string, err error) {
		log.
			ErrorMsg(msg).
			FieldErr(err).
			Write()
	})
}

func pollEventStream(projection taskprojection.TaskProjection, log logger.Log) {
	for {
		err := projection.CatchupEventStream()
		if err != nil {
			log.ErrorMsg("TaskProjection Polling").FieldErr(err).Write()
		}

		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
