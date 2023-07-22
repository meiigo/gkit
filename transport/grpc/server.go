package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/meiigo/gkit/log"
	"github.com/meiigo/gkit/transport"
	"github.com/meiigo/gkit/x/host"
)

var (
	_ transport.Server = (*Server)(nil)
)

// Server is a gRPC server wrapper
type Server struct {
	*grpc.Server
	lis       net.Listener
	tlsConf   *tls.Config
	once      sync.Once
	err       error
	network   string
	address   string
	endpoint  *url.URL
	timeout   time.Duration
	unaryInts []grpc.UnaryServerInterceptor
	grpcOpts  []grpc.ServerOption
	health    *health.Server
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":9090",
		timeout: 10 * time.Second,
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	// TODO: stream ints
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
	}

	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)

	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
// grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

// Start the gRPC server.
func (s *Server) Start(_ context.Context) error {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		s.lis = lis
	})
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	if s.err != nil {
		return s.err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop the gRPC server.
func (s *Server) Stop(_ context.Context) error {
	s.GracefulStop()
	s.health.Shutdown()
	log.Infof("[gRPC] server[%s] stopping", s.address)
	return nil
}
