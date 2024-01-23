package log

import "context"

type loggerKey struct{}

func FromContext(ctx context.Context) *Helper {
	l, ok := ctx.Value(loggerKey{}).(*Helper)
	if ok {
		return l
	}
	return NewHelper(DefaultLogger)
}

func NewContext(ctx context.Context, l *Helper) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}
