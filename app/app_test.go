package app_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/meiigo/gkit/app"
	"github.com/meiigo/gkit/monitor"
	transhttp "github.com/meiigo/gkit/transport/http"
)

// go test -v *.go -test.run=TestApp
func TestApp(t *testing.T) {
	hs := transhttp.NewServer(
		transhttp.Address(":8080"),
	)

	app := app.New(
		app.Name("gkit"),
		app.Version("v1.0.0"),
		app.Server(hs),
	)

	time.AfterFunc(60*time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestMux(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello Gorilla Mux!")
	})

	httpSrv := transhttp.NewServer(
		transhttp.Address(":8000"),
		transhttp.Handler(router),
	)

	app := app.New(
		app.Name("mux"),
		app.Version("v1.0.0"),
		app.Server(httpSrv),
	)
	time.AfterFunc(60*time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}

// curl -XGET http://127.0.0.1:8080/hello
// curl -XGET http://127.0.0.1:16060/monitor/env
func TestGin(t *testing.T) {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"hello": "world"})
	})

	httpSrv := transhttp.NewServer(
		transhttp.Address(":8080"),
		transhttp.Handler(router),
	)

	app := app.New(
		app.Name("gin"),
		app.Version("v1.0.0"),
		app.Server(httpSrv),
		app.Monitor(&monitor.Config{
			Enabled: true,
		}),
	)
	time.AfterFunc(60*time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
