package http

import (
	"fmt"

	"github.com/meiigo/gkit/examples/blog/internal/config"
	"github.com/meiigo/gkit/log"

	"github.com/gin-gonic/gin"
	"github.com/meiigo/gkit/app"
	"github.com/meiigo/gkit/examples/blog/internal/service"
	"github.com/meiigo/gkit/transport/http"
)

type Server struct {
	*http.Server
	blogSrv *service.BlogService
	log     log.Logger
}

func NewServer(c *app.HTTPServer, srv *service.BlogService) *Server {
	l := log.NewJSONLogger(log.WithLevel(config.GetLogLevel()))
	if config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.Use(gin.Recovery())
	opts := []http.ServerOption{
		http.Handler(engine),
	}
	if c.Port > 0 {
		opts = append(opts, http.Address(fmt.Sprintf(":%d", c.Port)))
	}
	s := &Server{
		log:     l,
		blogSrv: srv,
		Server:  http.NewServer(opts...),
	}
	s.Set(engine)
	return s
}
