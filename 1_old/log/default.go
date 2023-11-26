package log

import (
	"context"
)

var std ILog

func init() {
	std = NewDefaultLog()
}

func SetLogger(logger ILog) {
	std = logger
}

func Info(ctx context.Context, args ...interface{}) {
	std.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	std.Warn(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	std.Error(ctx, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	std.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	std.Warnf(ctx, format, args...)

}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	std.Errorf(ctx, format, args...)
}
