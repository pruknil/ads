package service

import (
	"bytes"
	"encoding/json"
	http2 "github.com/pruknil/ads/backends/http"
	"net/http"
)

type HttpBackendService struct {
	http2.IHttpBackendService
	http2.Config
}

func New(c http2.Config, service http2.IHttpBackendService) IHttpBackend {
	return &HttpBackendService{
		IHttpBackendService: service,
		Config:              c,
	}
}

func (s *HttpBackendService) newGetRequestBuilder(url string, query map[string]string, body interface{}) (http2.RestRequestBuilder, error) {
	var requestBuilder http2.RestRequestBuilder
	var reader *bytes.Reader
	if body != nil {
		requestBodyByte, _ := json.Marshal(body)
		reader = bytes.NewReader(requestBodyByte)
	}
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest("GET", url, reader)
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		return requestBuilder, err
	}
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	requestBuilder.Req = req
	return requestBuilder, nil
}

func (s *HttpBackendService) newPostRequestBuilder(url string, body interface{}) (http2.RestRequestBuilder, error) {
	var requestBuilder http2.RestRequestBuilder
	requestBodyByte, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBodyByte))
	if err != nil {
		return requestBuilder, err
	}
	requestBuilder.Req = req
	return requestBuilder, nil
}
