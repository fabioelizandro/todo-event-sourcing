package http_routes

import (
	"fabioelizandro/todo-event-sourcing/logger"
	"io/ioutil"
	"net/http"
)

type StdHttpRouteAdapter interface {
	Transform(Route) func(w http.ResponseWriter, r *http.Request)
}

type stdHttpRouteAdapter struct {
	log logger.Log
}

func (s *stdHttpRouteAdapter) Transform(route Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.
					ErrorMsg("Server: panicked").
					FieldInterface("panicked", err).
					Write()

				s.stdRouteAdapterSomethingWentWrong(w)
			}
		}()

		request, err := s.stdRouteAdapterRequest(r)
		if err != nil {
			panic(err)
		}

		response, err := route.Handle(request)
		if err != nil {
			s.log.
				ErrorMsg("Server: route handler error").
				FieldErr(err).
				Write()

			s.stdRouteAdapterSomethingWentWrong(w)
		}

		err = s.stdRouteAdapterResponse(response, w)
		if err != nil {
			s.log.
				ErrorMsg("Server: could not write response").
				FieldErr(err).
				Write()

			s.stdRouteAdapterSomethingWentWrong(w)
		}
	}
}

func (s *stdHttpRouteAdapter) stdRouteAdapterRequest(r *http.Request) (Request, error) {
	headers := Headers{}
	for key, _ := range r.Header {
		headers[key] = r.Header.Get(key)
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	requestBody, err := newRequestBody(headers, rawBody)
	if err != nil {
		return nil, err
	}

	return newRequest(headers, requestBody), nil
}

func (s *stdHttpRouteAdapter) stdRouteAdapterResponse(response Response, w http.ResponseWriter) error {
	body, err := response.Body()
	if err != nil {
		return err
	}

	for name, value := range response.Headers() {
		w.Header().Add(name, value)
	}

	err = w.Header().Write(w)
	if err != nil {
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (s *stdHttpRouteAdapter) stdRouteAdapterSomethingWentWrong(w http.ResponseWriter) {
	w.WriteHeader(500)
	_, err := w.Write([]byte("Something Went Wrong"))
	if err != nil {
		panic(err)
	}
}

func NewStdHttpRouteAdapter(log logger.Log) StdHttpRouteAdapter {
	return &stdHttpRouteAdapter{log}
}
