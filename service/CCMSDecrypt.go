package service

import (
	"encoding/json"
	"github.com/getlantern/deepcopy"
	"github.com/pruknil/ads/backends/http"
	"github.com/pruknil/ads/backends/http/service"
	"github.com/pruknil/ads/logger"
)

type CCMSDecryptService struct {
	baseService
	service.IHttpBackend
	logger.AppLog

	beRequest  http.DecryptDataRequest
	beResponse *http.DecryptDataResponse

	serviceRequest  CCMSDecryptRequest
	serviceResponse CCMSDecryptResponse
}

type CCMSDecryptRequest struct {
	DPKName string `json:"dPKName"`
	EData   string `json:"eData"`
}

type CCMSDecryptResponse struct {
	DPKName string `json:"dPKName"`
	Data    string `json:"data"`
}

func (s *CCMSDecryptService) Parse() error {
	jsonString, err := json.Marshal(s.Request.Body)
	if err != nil {
		return err
	}
	s.Log.Trace.Debug(jsonString)
	s.serviceRequest = CCMSDecryptRequest{}
	err = json.Unmarshal(jsonString, &s.serviceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *CCMSDecryptService) InputMapping() error {
	s.beRequest.KBankRequestHeader.UserId = s.Request.Header.UserId
	s.beRequest.KBankRequestHeader.RqDt = s.Request.Header.RqDt
	s.beRequest.KBankRequestHeader.FuncNm = "DecryptData"
	s.beRequest.KBankRequestHeader.RqAppId = "772"
	s.beRequest.KBankRequestHeader.RqUID = s.Request.Header.RqUID
	err := deepcopy.Copy(&s.beRequest.DecryptDataBodyRequest, s.serviceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *CCMSDecryptService) OutputMapping() error {
	//err := deepcopy.Copy(&s.serviceResponse, s.beResponse.DecryptDataBodyResponse)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (s *CCMSDecryptService) getResponse() ResMsg {
	s.Response.Body = s.serviceResponse
	return s.Response
}

func (s *CCMSDecryptService) Business() error {
	s.WSConn.WriteMessage(1,[]byte("Hello"))
	return nil
}
