package http_routes

import (
	"encoding/json"
	"strings"
)

type Headers map[string]string

func (h Headers) Merge(headers Headers) {
	for k, v := range headers {
		h[strings.ToLower(k)] = v
	}
}

type RequestBody interface {
	GetStr(field string, defaultValue string) string
}

type Request interface {
	Headers() Headers
	Body() RequestBody
}

type Response interface {
	Headers() Headers
	Body() ([]byte, error)
}

type request struct {
	headers Headers
	body    RequestBody
}

func (j *request) Headers() Headers {
	return j.headers
}

func (j *request) Body() RequestBody {
	return j.body
}

type jsonRequestBody struct {
	rawBody []byte
	body    map[string]string
}

func (j *jsonRequestBody) GetStr(field string, defaultValue string) string {
	value := j.body[field]
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

type jsonResponse struct {
	headers Headers
	payload interface{}
}

func (j *jsonResponse) Headers() Headers {
	return j.headers
}

func (j *jsonResponse) Body() ([]byte, error) {
	return json.Marshal(j.payload)
}

func newJsonResponse(headers Headers, payload interface{}) *jsonResponse {
	defaultHeaders := Headers{"status": "200", "content-type": "application/json"}
	defaultHeaders.Merge(headers)

	return &jsonResponse{headers: defaultHeaders, payload: payload}
}

func newJsonOkResponse(headers Headers) *jsonResponse {
	defaultHeaders := Headers{"status": "202"}
	defaultHeaders.Merge(headers)

	return newJsonResponse(defaultHeaders, map[string]string{"message": "OK"})
}

func NewRequest(headers Headers, body RequestBody) *request {
	return &request{
		headers: headers,
		body:    body,
	}
}

func NewJsonRequestBody(rawBody []byte) (*jsonRequestBody, error) {
	body := map[string]string{}
	err := json.Unmarshal(rawBody, body)
	if err != nil {
		return nil, err
	}

	return &jsonRequestBody{rawBody: rawBody, body: body}, nil
}
