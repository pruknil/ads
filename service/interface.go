package service

import (
	"context"
	"github.com/pruknil/ads/logger"
	"github.com/gorilla/websocket"
)

type IHttpService interface {
	DoService(ReqMsg, logger.AppLog, *websocket.Conn) ResMsg
}


type IServiceTemplate interface {
	Parse() error
	Validate() error
	OutputMapping() error
	InputMapping() error
	Business() error
	setRequest(ReqMsg) error
	getResponse() ResMsg
	DoService(req ReqMsg, service IServiceTemplate) (ResMsg, error)
	setLog(appLog logger.AppLog)
	setContext(context.Context)
	getContext() context.Context
	setWSConn(wsConn *websocket.Conn) error
}
