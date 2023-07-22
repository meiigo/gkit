// Monitor 对应用程序进行监控，通过restful api请求来监管、审计、收集应用的运行情况

package monitor

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/meiigo/gkit/monitor/env"
	"github.com/meiigo/gkit/monitor/metric"
	"github.com/meiigo/gkit/monitor/profile"
	transhttp "github.com/meiigo/gkit/transport/http"
)

type Server struct {
	*transhttp.Server
}

func New(conf Config) *Server {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	rg := engine.Group("/monitor")

	env.Router(rg)
	metric.Router(rg)
	profile.Route(rg)

	port := 16060
	if conf.Port > 0 {
		port = conf.Port
	}
	httpSrv := transhttp.NewServer(
		transhttp.Address(fmt.Sprintf(":%d", port)),
		transhttp.Handler(engine),
	)
	return &Server{
		Server: httpSrv,
	}
}

type Config struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
	Port    int  `json:"port" yaml:"port"`
}
