package http_essentials

import (
	"encoding/json"
	"strings"
)

type Route interface {
	Methods() []string
	Path() string
	Handle(Request) (Response, error)
}

type Headers map[string]string

func (h Headers) Merge(headers Headers) {
	for k, v := range headers {
		h[strings.ToLower(k)] = v
	}
}

func (h Headers) Value(name string, defaultValue string) string {
	value := h[strings.ToLower(name)]
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

type PathParams map[string]string

func (h PathParams) Value(name string, defaultValue string) string {
	value := h[strings.ToLower(name)]
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

type RequestBodyFields map[string]string

func (h RequestBodyFields) Value(name string, defaultValue string) string {
	value := h[strings.ToLower(name)]
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

type RequestBody interface {
	FieldStr(name string, defaultValue string) string
}

type Request interface {
	Headers() Headers
	Body() RequestBody
	PathParams() PathParams
}

type Response interface {
	Headers() Headers
	Body() ([]byte, error)
}

type request struct {
	headers    Headers
	body       RequestBody
	pathParams PathParams
}

func (j *request) Headers() Headers {
	return j.headers
}

func (j *request) Body() RequestBody {
	return j.body
}

func (j *request) PathParams() PathParams {
	return j.pathParams
}

type emptyRequestBody struct {
}

func (e *emptyRequestBody) FieldStr(field string, defaultValue string) string {
	return defaultValue
}

type unknownRequestBody struct {
}

func (e *unknownRequestBody) FieldStr(field string, defaultValue string) string {
	return defaultValue
}

type fakeRequestBody struct {
	fields RequestBodyFields
}

func (f *fakeRequestBody) FieldStr(name string, defaultValue string) string {
	return f.fields.Value(name, defaultValue)
}

type jsonRequestBody struct {
	rawBody []byte
	body    map[string]string
}

func (j *jsonRequestBody) FieldStr(name string, defaultValue string) string {
	value := j.body[name]
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

func NewJsonResponse(headers Headers, payload interface{}) *jsonResponse {
	defaultHeaders := Headers{"status": "200", "content-type": "application/json"}
	defaultHeaders.Merge(headers)

	return &jsonResponse{headers: defaultHeaders, payload: payload}
}

func NewJsonOkResponse(headers Headers) *jsonResponse {
	defaultHeaders := Headers{"status": "202"}
	defaultHeaders.Merge(headers)

	return NewJsonResponse(defaultHeaders, map[string]string{"message": "OK"})
}

func NewRequest(headers Headers, body RequestBody, pathParams PathParams) *request {
	return &request{
		headers:    headers,
		body:       body,
		pathParams: pathParams,
	}
}

func newRequestBody(headers Headers, rawBody []byte) (RequestBody, error) {
	if len(rawBody) == 0 {
		return newEmptyRequestBody(), nil
	}

	contentType := headers.Value("content-type", "")
	switch contentType {
	case "application/json":
		return newJsonRequestBody(rawBody)
	default:
		return newUnknownRequestBody(), nil
	}
}

func newJsonRequestBody(rawBody []byte) (*jsonRequestBody, error) {
	body := map[string]string{}
	err := json.Unmarshal(rawBody, body)
	if err != nil {
		return nil, err
	}

	return &jsonRequestBody{rawBody: rawBody, body: body}, nil
}

func newUnknownRequestBody() *unknownRequestBody {
	return &unknownRequestBody{}
}

func newEmptyRequestBody() *emptyRequestBody {
	return &emptyRequestBody{}
}

func NewFakeRequestBody(fields RequestBodyFields) *fakeRequestBody {
	return &fakeRequestBody{fields: fields}
}

func NewEmptyFakeRequestBody() *fakeRequestBody {
	return &fakeRequestBody{fields: RequestBodyFields{}}
}
