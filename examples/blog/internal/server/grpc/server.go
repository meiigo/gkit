package grpc

import (
	"fmt"

	"github.com/meiigo/gkit/app"
	"github.com/meiigo/gkit/examples/blog/api"
	tg "github.com/meiigo/gkit/transport/grpc"
	"google.golang.org/grpc"
)

func NewGRPCServer(c *app.GRPCServer, s api.BlogServer) *tg.Server {
	opts := []tg.ServerOption{
		tg.UnaryInterceptor([]grpc.UnaryServerInterceptor{
			//grpc_prometheus.UnaryServerInterceptor,
			//requestid.GRPCServerMiddleware(),
		}...),
	}
	if c.Port > 0 {
		opts = append(opts, tg.Address(fmt.Sprintf(":%d", c.Port)))
	}

	srv := tg.NewServer(opts...)
	api.RegisterBlogServer(srv, s)
	return srv
}
