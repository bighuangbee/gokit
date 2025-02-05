package log

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
)

// Logger interface
type Logger interface {
	// WithCtx wrap logger with context
	WithCtx(ctx context.Context) Logger
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Log(level klog.Level, keyvals ...interface{}) error
}
