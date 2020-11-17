package main

import (
	"github.com/pruknil/ads/app"
	httpbackend "github.com/pruknil/ads/backends/http"
	httpsvc "github.com/pruknil/ads/backends/http/service"
	"github.com/pruknil/ads/logger"
	"github.com/pruknil/ads/router"
	"github.com/pruknil/ads/router/http"
	"github.com/pruknil/ads/service"
	"go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	container := buildContainer()
	err := invokeContainer(container)
	if err != nil {
		log.Fatal("Invoke Container error")
	}
}

func invokeContainer(container *dig.Container) error {
	return container.Invoke(func(route []router.IRouter) {
		for _, v := range route {
			v.Start()
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-quit

		for _, v := range route {
			v.Shutdown()
		}
	})
}
func errorWrap(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
func buildContainer() *dig.Container {
	container := dig.New()
	errorWrap(container.Provide(NewConfig))
	errorWrap(container.Provide(NewLogger))

	errorWrap(container.Provide(NewHttpService))
	errorWrap(container.Provide(NewHttp))
	errorWrap(container.Provide(NewHttpBackend))

	//errorWrap(container.Provide(NewDataBase))
	//errorWrap(container.Provide(NewDataService))

	errorWrap(container.Provide(NewRouter))
	return container
}

func NewLogger() logger.AppLog {
	al := logger.New()
	al.Error = al.NewLog("error", os.Getenv("LOG_LEVEL"))
	al.Perf = al.NewLog("perf", os.Getenv("LOG_LEVEL"))
	al.Trace = al.NewLog("trace", os.Getenv("LOG_LEVEL"))
	al.Rest = al.NewLog("rest", os.Getenv("LOG_LEVEL"))
	al.Router = al.NewLog("router", os.Getenv("LOG_LEVEL"))
	al.Socket = al.NewLog("socket", os.Getenv("LOG_LEVEL"))
	return al
}

//================= Start BACKEND Section =================

func NewHttp(cfg app.Config, log logger.AppLog) httpbackend.IHttpBackendService {
	return httpbackend.New(cfg.Backend.Http, log)
}

func NewHttpBackend(cfg app.Config, s httpbackend.IHttpBackendService) httpsvc.IHttpBackend {
	return httpsvc.New(cfg.Backend.Http, s)
}

//================= End BACKEND Section =================

//Create all router here eg.. rest, socket, sftp
func NewRouter(httpService service.IHttpService, conf app.Config, log logger.AppLog) []router.IRouter {
	var route []router.IRouter
	route = append(route, http.NewGin(conf.Router.Http, httpService, log))
	return route
}

// service
//Http service
func NewHttpService(httpService httpsvc.IHttpBackend, cfg app.Config, log logger.AppLog) service.IHttpService {
	return service.HttpService{
		//CCMSDecryptService: service.CCMSDecryptService{IHttpBackend: httpService, AppLog: log},
		//CCMSEncryptService: service.CCMSEncryptService{IHttpBackend: httpService, AppLog: log},
		//DemoService:        service.DemoService{AppLog: log, CcmsFs: syn},
		//SFTP0001O01Service: service.SFTP0001O01Service{AppLog: log, CcmsFs: syn},
		//SFTP0002I01Service: service.SFTP0002I01Service{AppLog: log, CcmsFs: syn},
	}
}
func NewConfig() app.Config {
	ten, _ := time.ParseDuration("10s")
	return app.Config{
		Backend: app.Backend{
			Http: httpbackend.Config{
				CCMSAPI: httpbackend.CCMSAPI{
					DecryptDataUrl: "https://172.30.70.51:9000/v2/DecryptData",
					EncryptDataUrl: "https://172.30.70.51:9000/v2/EncryptData",
					LocalUrl:       "http://localhost:8080/api",
				},
			},
		},
		Router: app.Router{
			Http: http.Config{
				Port:         os.Getenv("SERVER_PORT"),
				ReadTimeout:  ten,
				WriteTimeout: ten,
				IdleTimeout:  ten,
			},
		},
		Service: app.Service{
			Http: service.Config{},
		},
	}
}
