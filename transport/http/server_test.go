package http_test

import (
	"context"
	"testing"
	"time"

	transhttp "github.com/meiigo/gkit/transport/http"
)

func TestServer_Start(t *testing.T) {
	srv := transhttp.NewServer(
		transhttp.Address(":8080"),
	)

	go func() {
		if err := srv.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	time.AfterFunc(60*time.Second, func() {
		defer srv.Stop(context.TODO())

	})
	time.Sleep(100 * time.Second)
}
