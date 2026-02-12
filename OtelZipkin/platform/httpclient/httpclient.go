package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type a struct {
}

func New() HTTPClient {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func (h *a) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
