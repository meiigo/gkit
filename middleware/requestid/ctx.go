package requestid

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	key = "x-request-id"
)

type requestIDKey struct {
}

func NewContext(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, v)
}

func FromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(requestIDKey{}).(string)
	return val, ok
}

func FromIncoming(ctx context.Context) (string, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(key)
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return "", false
}

func ToOutgoing(ctx context.Context, v string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, v)
}
