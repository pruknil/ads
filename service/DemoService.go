package service

import (
	"encoding/json"
	"github.com/getlantern/deepcopy"
	"github.com/pruknil/ads/backends/http"
	"github.com/pruknil/ads/backends/http/service"
	"github.com/pruknil/ads/logger"
)

type DemoService struct {
	baseService
	service.IHttpBackend
	logger.AppLog

	beRequest  http.DecryptDataRequest
	beResponse *http.DecryptDataResponse

	serviceRequest  DemoRequest
	serviceResponse DemoResponse
}

type DemoRequest struct {
	DPKName string `json:"dPKName"`
	EData   string `json:"eData"`
}

type DemoResponse struct {
	DPKName string `json:"dPKName"`
	Data    string `json:"data"`
}

func (s *DemoService) Parse() error {
	jsonString, err := json.Marshal(s.Request.Body)
	if err != nil {
		return err
	}
	s.Log.Trace.Debug(jsonString)
	s.serviceRequest = DemoRequest{}
	err = json.Unmarshal(jsonString, &s.serviceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *DemoService) InputMapping() error {
	s.beRequest.RequestHeader.UserId = s.Request.Header.UserId
	s.beRequest.RequestHeader.RqDt = s.Request.Header.RqDt
	s.beRequest.RequestHeader.RqUID = s.Request.Header.RqUID
	err := deepcopy.Copy(&s.beRequest.DecryptDataBodyRequest, s.serviceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *DemoService) OutputMapping() error {
	//err := deepcopy.Copy(&s.serviceResponse, s.beResponse.DecryptDataBodyResponse)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (s *DemoService) getResponse() ResMsg {
	s.Response.Body = s.serviceResponse
	return s.Response
}

func (s *DemoService) Business() error {
	s.WSConn.WriteMessage(1, []byte("2.jpg"))
	return nil
}
