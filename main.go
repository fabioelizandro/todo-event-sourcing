package main

import (
	"fabioelizandro/todo-event-sourcing/eventstream"
	"fabioelizandro/todo-event-sourcing/http_routes"
	"fabioelizandro/todo-event-sourcing/task"
	"fabioelizandro/todo-event-sourcing/taskprojection"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	stream := eventstream.NewInMemoryEventStream()
	commandHandler := task.NewCmdHandler(stream)
	projection := taskprojection.NewTaskProjection(stream)
	routes := []http_routes.Route{
		http_routes.NewTaskListRoute(projection),
		http_routes.NewTaskCreateRoute(commandHandler),
	}

	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.Path(), stdHttpRouteAdapter(route)).Methods(route.Methods()...)
	}

	go projection.PollEventStream(100)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func stdHttpRouteAdapter(route http_routes.Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Refactor and make it detect the request type instead of assume is always json

		headers := http_routes.Headers{}
		for key, _ := range r.Header {
			headers[key] = r.Header.Get(key)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		requestBody, err := http_routes.NewJsonRequestBody(body)
		if err != nil {
			panic(err)
		}

		request := http_routes.NewRequest(headers, requestBody)
		response, err := route.Handle(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Something went wrong"))
			panic(err)
		} else {
			body, err := response.Body()
			if err != nil {
				panic(err)
			}
			status, err := strconv.Atoi(response.Headers()["status"])
			if err != nil {
				panic(err)
			}

			w.WriteHeader(status)
			_, err = w.Write(body)
			if err != nil {
				panic(err)
			}
		}
	}
}
