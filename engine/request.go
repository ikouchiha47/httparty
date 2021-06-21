package engine

import (
	"bytes"
	"net/http"
	"strings"
)

type Body map[string]interface{}
type Headers map[string]string

type Request struct {
	URL     string  `json:"url"`
	Method  string  `json:"method,omitempty"`
	Type    string  `json:"type,omitempty"`
	Accept  string  `json:"accept,omitempty"`
	Headers Headers `json:"headers,omitempty"`
	Body    Body    `json:"body,omitempty"`
}

func (r Request) ReqMethod() string {
	method := strings.ToLower(r.Method)
	if method == "" {
		return "get"
	}

	return r.Method
}

func (r Request) ReqType() string {
	typ := strings.ToLower(r.Type)
	if typ == "" {
		return JSON
	}

	return typ
}

func (r Request) AcceptType() string {
	typ := strings.ToLower(r.Accept)
	if typ == "" {
		return JSON
	}

	return typ
}

func (r Request) ShouldMakeBody() bool {
	m := r.ReqMethod()

	return m == "post" || m == "put" || m == "patch"
}

func (r Request) MakeGet(body string) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = body

	return req, nil
}

func (r Request) MakePost(body string) (*http.Request, error) {
	if body == "" {
		return http.NewRequest(r.Method, r.URL, nil)
	}

	return http.NewRequest(r.Method, r.URL, bytes.NewBufferString(body))
}

type Response struct {
	Body Body
}
