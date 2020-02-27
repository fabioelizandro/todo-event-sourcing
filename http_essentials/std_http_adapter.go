package http_essentials

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type StdHttpRouteAdapter interface {
	Transform(Route) func(w http.ResponseWriter, r *http.Request)
}

type stdHttpRouteAdapter struct {
	onError func(string, error)
}

func (s *stdHttpRouteAdapter) Transform(route Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.onError("Server: panicked", err.(error))
				s.stdRouteAdapterSomethingWentWrong(w)
			}
		}()

		request, err := s.stdRouteAdapterRequest(r)
		if err != nil {
			panic(err)
		}

		response, err := route.Handle(request)
		if err != nil {
			s.onError("Server: route handler error", err)
			s.stdRouteAdapterSomethingWentWrong(w)
		}

		err = s.stdRouteAdapterResponse(response, w)
		if err != nil {
			s.onError("Server: could not write response", err)
			s.stdRouteAdapterSomethingWentWrong(w)
		}
	}
}

func (s *stdHttpRouteAdapter) stdRouteAdapterRequest(r *http.Request) (Request, error) {
	headers := Headers{}
	for key, _ := range r.Header {
		headers[strings.ToLower(key)] = r.Header.Get(key)
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	requestBody, err := newRequestBody(headers, rawBody)
	if err != nil {
		return nil, err
	}

	pathParams := PathParams{}
	for key, value := range mux.Vars(r) {
		pathParams[strings.ToLower(key)] = value
	}

	return NewRequest(headers, requestBody, pathParams), nil
}

func (s *stdHttpRouteAdapter) stdRouteAdapterResponse(response Response, w http.ResponseWriter) error {
	body, err := response.Body()
	if err != nil {
		return err
	}

	for name, value := range response.Headers() {
		w.Header().Set(name, value)
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

func NewStdHttpRouteAdapter(onError func(string, error)) StdHttpRouteAdapter {
	return &stdHttpRouteAdapter{onError: onError}
}
