package app

import (
	"context"
	"os"
	"time"

	"github.com/meiigo/gkit/log"
	"github.com/meiigo/gkit/monitor"

	"github.com/meiigo/gkit/transport"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id       string
	name     string
	version  string
	metadata map[string]string

	ctx  context.Context
	sigs []os.Signal

	logger      *log.Helper
	stopTimeout time.Duration
	servers     []transport.Server
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// StopTimeout with app stop timeout.
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// Logger with logger.
func Logger(l log.Logger) Option {
	return func(o *options) { o.logger = log.NewHelper(l) }
}

// Sigs with Sigs.
func Sigs(sigs []os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Context with context.
func Context(c context.Context) Option {
	return func(o *options) { o.ctx = c }
}

// Monitor with app monitor.
func Monitor(conf *monitor.Config) Option {
	if conf != nil && conf.Enabled {
		return func(o *options) {
			o.servers = append(o.servers, monitor.New(*conf))
		}
	}
	return func(o *options) {}
}
