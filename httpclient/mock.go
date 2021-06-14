package httpclient

import "net/http"

type ClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}

	return &http.Response{}, nil
}

