package httpclient

import (
	"net/http"
	"time"
)

type HttpClientI interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	*http.Client
}

func NewHttpClient(timeoutInSec int64) *HttpClient {
	return &HttpClient{
		&http.Client{
			Timeout: time.Second * time.Duration(timeoutInSec),
		},
	}
}
