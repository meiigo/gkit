package log_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/meiigo/gkit/log"
)

// go test -v *.go -test.run=^TestLogger$
func TestLogger(t *testing.T) {

	fmt.Println("log...")
	log.Info("this is a msg")
	log.Infof("this is a %s msg", "format")

	// default log
	fmt.Println("default log...")
	dl := log.NewLogger(log.WithLevel(log.LevelInfo))
	dl.Log(log.LevelDebug, "msg", "this is a msg")
	dl.Log(log.LevelInfo, "msg", "this is a msg")

	// json log
	fmt.Println("json log...")
	jl := log.NewJsonLogger(log.WithLevel(log.LevelInfo), log.WithOutput(os.Stdout))
	jl.Log(log.LevelDebug, "msg", "this is a msg")
	jl.Log(log.LevelInfo, "msg", "this is a msg")

	// helper
	fmt.Println("helper...")
	h1 := log.NewHelper(log.NewLogger(log.WithLevel(log.LevelInfo)))
	h1.Info("this is a msg")
	h1.Infof("this is a %s msg", "format")
	h1.Infow("a1", "b1", "a2", "b2")

	h2 := h1.WithFields(map[string]interface{}{"k1": "v1"})
	h2.Info("this is a msg")
	h2.Infof("this is a %s msg", "format")
	h2.Infow("a1", "b1", "a2", "b2")

	h3 := h2.WithFields(map[string]interface{}{"k2": "v2"})
	h3.Info("this is a msg")
	h3.Infof("this is a %s msg", "format")
	h3.Infow("a1", "b1", "a2", "b2")

	// error
	fmt.Println("error...")
	err := errors.New("this is a error")
	h1.WithError(err).Info()
	h3.WithError(err).Info()
}
