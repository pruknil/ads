package service

import (
	"encoding/json"
	http2 "github.com/pruknil/ads/backends/http"
)

func (s *HttpBackendService) DecryptData(req http2.DecryptDataRequest) (*http2.DecryptDataResponse, error) {
	var model http2.DecryptDataResponse
	var builder, err = s.newPostRequestBuilder(s.CCMSAPI.DecryptDataUrl, req)
	if err != nil {
		return nil, err
	}
	builder.
		SetResponseModel(&model).
		SetSecurityTransport().
		SetContentType("application/json")

	err = s.IHttpBackendService.DoRequest(builder)
	if err != nil {
		errX := json.Unmarshal([]byte(err.Error()), &model)
		if errX == nil {
			return nil, err
		}
		return &model, err
	}
	return &model, nil
}

func (s *HttpBackendService) EncryptData(req http2.EncryptDataRequest) (*http2.EncryptDataResponse, error) {
	var model http2.EncryptDataResponse
	var builder, err = s.newPostRequestBuilder(s.CCMSAPI.EncryptDataUrl, req)
	if err != nil {
		return nil, err
	}
	builder.
		SetResponseModel(&model).
		SetSecurityTransport().
		SetContentType("application/json")

	err = s.IHttpBackendService.DoRequest(builder)
	if err != nil {
		errX := json.Unmarshal([]byte(err.Error()), &model)
		if errX == nil {
			return nil, err
		}
		return &model, err
	}
	return &model, nil
}

func (s *HttpBackendService) SFTP0002I01(req http2.SFTP0002I01Request) (*http2.SFTP0002I01Response, error) {
	var model http2.SFTP0002I01Response
	var builder, err = s.newPostRequestBuilder(s.CCMSAPI.LocalUrl, req)
	if err != nil {
		return nil, err
	}
	builder.
		SetResponseModel(&model).
		SetSecurityTransport().
		SetContentType("application/json")

	err = s.IHttpBackendService.DoRequest(builder)
	if err != nil {
		errX := json.Unmarshal([]byte(err.Error()), &model)
		if errX == nil {
			return nil, err
		}
		return &model, err
	}
	return &model, nil
}
