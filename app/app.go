package app

import (
	http2 "github.com/pruknil/ads/backends/http"
	"github.com/pruknil/ads/logger"
	"github.com/pruknil/ads/router/http"
	"github.com/pruknil/ads/service"
)

type Config struct {
	logger.AppLog
	Backend
	Router
	Service
}

type Router struct {
	Http http.Config
}

type Service struct {
	Http service.Config
}

type Backend struct {
	Http http2.Config
}
