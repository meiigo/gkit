package main

import (
	"flag"

	"github.com/meiigo/gkit/app"
	"github.com/meiigo/gkit/examples/blog/internal/config"
	"github.com/meiigo/gkit/examples/blog/internal/server/grpc"
	"github.com/meiigo/gkit/examples/blog/internal/server/http"
	"github.com/meiigo/gkit/examples/blog/internal/service"
	"github.com/meiigo/gkit/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {

	flag.Parse()
	log.Infof("start server...")

	if err := config.InitConfig(flagconf); err != nil {
		panic(err)
	}

	blogSrv := service.NewBlogService(nil, log.DefaultLogger)

	httpSrv := http.NewServer(config.GetHttpConfig(), blogSrv)
	grpcSrv := grpc.NewServer(config.GetGRPCConfig(), blogSrv)

	a := app.New(
		app.ID("app-server"),
		app.Name("ab-server"),
		app.Version("v1.0.0"),
		app.Server(httpSrv, grpcSrv),
		app.Monitor(config.GetMonitorConfig()),
	)

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}

	log.Infof("stop server...")
}
