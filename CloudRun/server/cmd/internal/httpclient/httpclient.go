package httpclient

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type a struct {
	
}

func New() HTTPClient{
	return &a{}
}

func (h *a) Do(req *http.Request) (*http.Response, error){
	return http.DefaultClient.Do(req)
}
