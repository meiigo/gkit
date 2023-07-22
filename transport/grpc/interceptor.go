package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/meiigo/gkit/transport"
)

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var cancel context.CancelFunc
		md, _ := metadata.FromIncomingContext(ctx)
		replyHeader := metadata.MD{}
		ctx = transport.NewServerContext(ctx, &Transport{
			endpoint:    s.endpoint.String(),
			operation:   info.FullMethod,
			reqHeader:   grpcHeader(md),
			replyHeader: grpcHeader(replyHeader),
		})
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

//func unaryClientInterceptor(ms []endpoint.Middleware, timeout time.Duration) grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//		ctx = transport.NewClientContext(ctx, &Transport{
//			endpoint:  cc.Target(),
//			operation: method,
//			reqHeader: grpcHeader{},
//		})
//		if timeout > 0 {
//			var cancel context.CancelFunc
//			ctx, cancel = context.WithTimeout(ctx, timeout)
//			defer cancel()
//		}
//		h := func(ctx context.Context, req interface{}) (interface{}, error) {
//			if tr, ok := transport.FromClientContext(ctx); ok {
//				header := tr.RequestHeader()
//				keys := header.Keys()
//				kvs := make([]string, 0, len(keys))
//				for _, k := range keys {
//					kvs = append(kvs, k, header.Get(k))
//				}
//				ctx = metadata.AppendToOutgoingContext(ctx, kvs...)
//			}
//			return reply, invoker(ctx, method, req, reply, cc, opts...)
//		}
//		if len(ms) > 0 {
//			h = endpoint.Chain(ms...)(h)
//		}
//		_, err := h(ctx, req)
//		return err
//	}
//}
