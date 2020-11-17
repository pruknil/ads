package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pruknil/ads/logger"
	"github.com/pruknil/ads/service"
	"go.uber.org/dig"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Gin struct {
	srv *http.Server
	//srvSSL      *http.Server
	config      Config
	httpService service.IHttpService
	router      *gin.Engine
	log         logger.AppLog
	container   *dig.Container
	wsConn      *websocket.Conn
}

func NewGin(cfg Config, service service.IHttpService, logg logger.AppLog) *Gin {
	return &Gin{
		config:      cfg,
		httpService: service,
		log:         logg,
	}
}

func (g *Gin) initializeRoutes() {
	hn, _ := os.Hostname()
	//g.router.Use(static.Serve("/", static.LocalFile("views/static", true)))
	g.router.POST("/api", g.serviceLocator)

	g.router.GET("/health", func(c *gin.Context) {
		c.String(200, hn)

	})
	g.router.GET("/ws", g.atmMachine)
}

var wsupgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  nil,
	Subprotocols:     nil,
	Error:            nil,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

func (g *Gin) atmMachine(c *gin.Context) {
	//g.wshandler(c.Writer, c.Request)
	var err error
	g.wsConn, err = wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	for {
		t, msg, err := g.wsConn.ReadMessage()
		if err != nil {
			break
		}
		g.wsConn.WriteMessage(t, msg)
		log.Println(msg)
	}
}

func (g *Gin) serviceLocator(c *gin.Context) {
	var reqMsg service.ReqMsg
	err := c.BindJSON(&reqMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resMsg := g.httpService.DoService(reqMsg, g.log, g.wsConn)
	c.JSON(http.StatusOK, resMsg)

}

func (g *Gin) Start() {
	//g.container = container
	g.router = gin.Default()

	g.router.Use(GenRsUID())
	g.router.Use(LogRequest(&g.log))
	g.router.Use(LogResponse(&g.log))

	g.initializeRoutes()

	go func() {
		g.srv = &http.Server{
			Addr:         ":" + g.config.Port,
			Handler:      g.router,
			ReadTimeout:  g.config.ReadTimeout,
			WriteTimeout: g.config.WriteTimeout,
			IdleTimeout:  g.config.IdleTimeout,
		}
		if err := g.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (g *Gin) Shutdown() {
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
