package grpc

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/meiigo/gkit/log"
	"github.com/meiigo/gkit/middleware/circuitbreaker"
	"google.golang.org/grpc"
	secure "google.golang.org/grpc/credentials/insecure"
)

// ClientOption is gRPC client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with client endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithOptions with gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint string
	timeout  time.Duration
	grpcOpts []grpc.DialOption
	logger   log.Logger
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

// Dial returns a GRPC connection.
func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	l := log.NewHelper(log.DefaultLogger)
	options := clientOptions{
		timeout: 3000 * time.Millisecond, // 默认3秒超时
		logger:  l,
	}
	for _, o := range opts {
		o(&options)
	}

	ints := []grpc.UnaryClientInterceptor{
		circuitbreaker.GRPCClientCircuitBreaker(hystrix.CommandConfig{
			Timeout:                10,
			MaxConcurrentRequests:  10000,
			RequestVolumeThreshold: 10,
			SleepWindow:            10,
			ErrorPercentThreshold:  30,
		}),
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(ints...),
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(secure.NewCredentials()))
	}

	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
