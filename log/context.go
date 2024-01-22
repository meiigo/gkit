package log

import "context"

type loggerKey struct{}

func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey{}).(Logger)
	if ok {
		return l
	}
	return DefaultLogger
}

func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}
