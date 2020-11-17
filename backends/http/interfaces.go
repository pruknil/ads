package http

type IHttpBackendService interface {
	DoRequest(RestRequestBuilder) error
}
