package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/meiigo/gkit/log"
	"github.com/meiigo/gkit/x/host"
	"github.com/pkg/errors"
)

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL
	network  string
	address  string
	timeout  time.Duration
	handler  http.Handler
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
	}
	for _, o := range opts {
		o(s)
	}
	s.Server = &http.Server{
		Handler:   s.handler,
		TLSConfig: s.tlsConf,
	}
	return s
}

// Start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{
		Scheme: "http",
		Host:   addr,
	}
	log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Infof("[HTTP] server[%s] stopping", s.Addr)
	return s.Shutdown(ctx)
}
