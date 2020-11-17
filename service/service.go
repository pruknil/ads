package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/pruknil/ads/logger"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

type Config struct {
}
type baseService struct {
	context.Context
	IServiceTemplate
	Request  ReqMsg
	Response ResMsg
	WSConn   *websocket.Conn
	Log      logger.AppLog
}

func (s *baseService) buildResponseHeader() ResHeader {
	return ResHeader{
		FuncNm:     s.Request.Header.FuncNm,
		RqUID:      s.Request.Header.RqUID,
		RsAppID:    ApplicationId,
		RsUID:      s.Request.Header.RsUID,
		RsDt:       time.Now(),
		StatusCode: "00",
		ErrorVect:  nil,
	}
}
func (s *baseService) getResponse() ResMsg {
	s.Response.Header = s.buildResponseHeader()
	return s.Response
}

func (s *baseService) setRequest(r ReqMsg) error {
	r.Header.RsUID = uuid.New().String()
	s.Request = r
	s.Context = context.WithValue(s.getContext(), "GUID", r.Header.RsUID)
	s.Response = ResMsg{}
	return nil
}

func (s *baseService) setContext(ctx context.Context) {
	s.Context = ctx
}

func (s *baseService) getContext() context.Context {
	return s.Context
}
func (s *baseService) setWSConn(wsConn *websocket.Conn) error {
	s.WSConn = wsConn
	return nil
}
func (s *baseService) callService(funcNm string, msg ReqMsg, body interface{}) (ResMsg, error) {
	servicesCtx := s.getContext().Value("SERVICES")
	serviceMap := servicesCtx.(HttpService)
	service := serviceMap.collectService(funcNm)
	service.setLog(s.Log)
	service.setContext(s.Context)
	h := ReqHeader{}
	err := deepcopy.Copy(&h, msg.Header)
	if err != nil {
		return ResMsg{}, err
	}
	h.FuncNm = funcNm
	return service.DoService(ReqMsg{
		Header: h,
		Body:   body,
	}, service)
}
func (s *baseService) setLog(appLog logger.AppLog) {
	s.Log = appLog
}

func (s *baseService) Parse() error {
	return nil
}
func (s *baseService) Validate() error {
	return nil
}
func (s *baseService) InputMapping() error {
	return nil
}
func (s *baseService) DoService(req ReqMsg, service IServiceTemplate) (ResMsg, error) {
	resMsg := ResMsg{}
	s.Log.Perf.WithFields(logrus.Fields{
		logger.FUNCNM: req.Header.FuncNm,
		logger.RQUID:  req.Header.RqUID,
		logger.RSUID:  "",
		logger.STATUS: "",
		logger.TM:     "",
	}).Print("I")
	defer func(t time.Time, sv *baseService, msg *ResMsg) {
		sv.Log.Perf.WithFields(logrus.Fields{
			logger.FUNCNM: msg.Header.FuncNm,
			logger.RQUID:  msg.Header.RqUID,
			logger.RSUID:  msg.Header.RsUID,
			logger.STATUS: msg.Header.StatusCode,
			logger.TM:     fmt.Sprintf("%.4f", float64(time.Since(t).Nanoseconds())/float64(time.Millisecond)),
		}).Print("O")
	}(time.Now(), s, &resMsg)
	setRequestErr := service.setRequest(req)
	if setRequestErr != nil {
		resMsg = s.buildErrResp(setRequestErr)
		return resMsg, errors.New("parse error")
	}
	parseErr := service.Parse()
	if parseErr != nil {
		resMsg = s.buildErrResp(parseErr)
		return resMsg, errors.New("parse error")
	}
	validateErr := service.Validate()
	if validateErr != nil {
		resMsg = s.buildErrResp(validateErr)
		return resMsg, errors.New("validate error")
	}

	inputMappingErr := service.InputMapping()
	if inputMappingErr != nil {
		resMsg = s.buildErrResp(inputMappingErr)
		return resMsg, errors.New("InputMapping Error")
	}

	businessError := service.Business()
	if businessError != nil {
		resMsg = s.buildErrResp(businessError)
		return resMsg, businessError
	}
	outputMappingErr := service.OutputMapping()
	if outputMappingErr != nil {
		resMsg = s.buildErrResp(outputMappingErr)
		return resMsg, errors.New("OutputMapping Error")
	}
	resMsg = service.getResponse()
	return resMsg, nil
}
func (s *baseService) buildErrResp(businessError error) ResMsg {
	return ResMsg{
		Header: ResHeader{
			FuncNm:     s.Request.Header.FuncNm,
			RqUID:      s.Request.Header.RqUID,
			RsUID:      s.Request.Header.RsUID,
			RsDt:       time.Now(),
			RsAppID:    ApplicationId,
			StatusCode: "10",
			ErrorVect:  &ErrorVect{Error: []Error{buildError(businessError)}},
		},
		Body: s.Response.Body,
	}
}
func buildError(businessError error) Error {
	var code, desc string
	sps := strings.Split(businessError.Error(), ":")
	if len(sps) > 1 {
		code = sps[0]
		desc = sps[1]
	} else {
		desc = businessError.Error()
	}
	return Error{
		ErrorAppID:    ApplicationId,
		ErrorAppAbbrv: ApplicationABBR,
		ErrorCode:     code,
		ErrorDesc:     desc,
		ErrorSeverity: "00",
	}
}

type HttpService struct {
	baseService
	CCMSDecryptService
	//DecipherAesService
	//EncipherAesService
	//DemoService
	//CCMSDecryptService
	//CCMSEncryptService
	//SFTP0001O01Service
	//SFTP0002I01Service
}

func (s *HttpService) collectService(serviceId string) IServiceTemplate {
	switch serviceId {
	//case "DecipherAes":
	//	return &s.DecipherAesService
	//case "EncipherAes":
	//	return &s.EncipherAesService
	//case "Demo":
	//	return &s.DemoService
	case "CCMSDecrypt":
		return &s.CCMSDecryptService
	//case "CCMSEncrypt":
	//	return &s.CCMSEncryptService
	//case "SFTP0001O01":
	//	return &s.SFTP0001O01Service
	//case "SFTP0002I01":
	//	return &s.SFTP0002I01Service
	}
	return nil
}
func (s HttpService) DoService(req ReqMsg, l logger.AppLog, wsConn *websocket.Conn) ResMsg {
	sv := s.collectService(req.Header.FuncNm)
	if sv == nil {
		return ResMsg{
			Header: ResHeader{
				RsDt:       time.Now(),
				StatusCode: "10",
				ErrorVect: &ErrorVect{[]Error{{
					ErrorAppID:    ApplicationId,
					ErrorAppAbbrv: ApplicationABBR,
					ErrorDesc:     "Service Not found",
					ErrorSeverity: "00",
				}}},
			},
			Body: nil,
		}
	}
	sv.setLog(l)
	if s.getContext() == nil {
		ctx := context.WithValue(context.Background(), "SERVICES", s)
		sv.setContext(ctx)
	}
	sv.setWSConn(wsConn)
	r, _ := sv.DoService(req, sv)
	return r
}
