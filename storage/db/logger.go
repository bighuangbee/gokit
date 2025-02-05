package db

import (
	"context"
	"errors"
	kitLog "github.com/bighuangbee/gokit/log"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"time"
)

func NewLogger(logger *kitLog.ZapLogger) *Logger {
	return &Logger{L: log.NewHelper(logger)}
}

// Logger Logger for gorm2
type Logger struct {
	L *log.Helper
}

// LogMode LogMode
func (l Logger) LogMode(glog.LogLevel) glog.Interface {
	return l
}

// Info Info
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.L.WithContext(ctx).Debugf(msg, args...)
}

// Warn Warn
func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.L.WithContext(ctx).Infof(msg, args...)
}

// Error Error
func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.L.WithContext(ctx).Errorf(msg, args...)
}

// Trace Trace, notfound和uniqueKey的错误不会印error
// TODO: 慢查询
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	spend := elapsed.Milliseconds()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			goto normalLog
		}
		if is, _, _ := IsUniqueErr(err); is {
			goto normalLog
		}
		l.L.WithContext(ctx).Errorw(
			"action", "sqlexec",
			"sql", sql,
			"spend", spend,
			"rows", rows,
			"err", err,
		)
		return
	}
normalLog:
	l.L.WithContext(ctx).Debugw("action", "sqlexec",
		"sql", sql,
		"spend", spend,
		"rows", rows,
	)
}
