package log

import (
	"context"
)

type logCtxKey struct{}

func Context(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, logger)
}

func getLogger(ctx context.Context) Logger {
	l, ok := ctx.Value(logCtxKey{}).(Logger)
	if ok {
		return l
	}
	return DefaultLogger
}

func With(ctx context.Context, fields ...Field) context.Context {
	logger := getLogger(ctx).With(fields...)
	return context.WithValue(ctx, logCtxKey{}, logger)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Warn(msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Panic(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Fatal(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Error(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	getLogger(ctx).Debug(msg, fields...)
}
