package circuitbreaker

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/meiigo/gkit/log"
	"google.golang.org/grpc"
)

// GRPCClientCircuitBreaker is a circuit breaker for grpc client middleware.
func GRPCClientCircuitBreaker(conf hystrix.CommandConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		m := hystrix.GetCircuitSettings()
		if _, ok := m[method]; !ok {
			hystrix.ConfigureCommand(method, conf)
		}
		err := hystrix.Do(method, func() (err error) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}, nil)
		cbs, _, _ := hystrix.GetCircuit(method)
		log.Debugf("用时:%s, 熔断器状态:%v, 请求是否允许:%v, error:%v\n", time.Since(start).String(), cbs.IsOpen(), cbs.AllowRequest(), err)
		return err
	}
}
