package zapx

import (
	"context"

	"go.uber.org/zap"
)

type zapLoggerKey struct{}

// Ctx returns a logger from ctx. If no logger is found in ctx, it returns
// the global logger.
func Ctx(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(zapLoggerKey{}).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return l
}

// With returns a copy of parent with the given logger.
func With(parent context.Context, l *zap.Logger) context.Context {
	return context.WithValue(parent, zapLoggerKey{}, l)
}

// WithFields returns a copy of parent with a logger created by adding fs to
// the logger from parent.
func WithFields(parent context.Context, fs ...zap.Field) context.Context {
	if len(fs) == 0 {
		panic("no fields provided")
	}
	return With(parent, Ctx(parent).With(fs...))
}
