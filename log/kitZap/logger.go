package kitZap

import (
	"context"
	"fmt"

	"git.hiscene.net/hifoundry/go-kit/constance/hiCommon"
	"git.hiscene.net/hifoundry/go-kit/log"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options options
type Options struct {
	// Level level, default INFO
	Level zapcore.Level
	// ServiceName serviceName
	ServiceName string
	// Skip caller skip, default 0
	Skip              int
	DisableStacktrace bool
}

type ZapLogger struct {
	*zap.SugaredLogger
}

func New(opt *Options) *ZapLogger {
	if opt == nil {
		panic("Options is required")
	}
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(opt.Level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "zmsg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		DisableStacktrace: opt.DisableStacktrace,
		DisableCaller:     false,
	}
	zapopts := []zap.Option{}
	if opt.Skip != 0 {
		zapopts = append(zapopts, zap.AddCallerSkip(opt.Skip))
	}
	l, err := cfg.Build(zapopts...)
	if err != nil {
		panic(err)
	}
	if opt.ServiceName != "" {
		l = l.With(zap.String("service", opt.ServiceName))
	}
	return &ZapLogger{
		SugaredLogger: l.Sugar(),
	}
}

func (l *ZapLogger) WithCtx(ctx context.Context) log.Logger {
	r := &ZapLogger{}

	v := tracing.TraceID()
	tid := v(ctx).(string)
	if tid != "" {
		r.SugaredLogger = l.SugaredLogger.With(zap.String(hiCommon.TraceID, v(ctx).(string)))
	} else {
		r.SugaredLogger = l.SugaredLogger
	}
	return r
}

// AddCallerSkip 继承配置并增加callerskip，返回一个新的ZapLogger
func (l *ZapLogger) AddCallerSkip(skip int) *ZapLogger {
	return &ZapLogger{
		SugaredLogger: l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar(),
	}
}

// Log for hiKratos
// see https://go-kratos.dev/docs/component/log
func (l *ZapLogger) Log(level klog.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	zl := l.SugaredLogger.Desugar()
	switch level {
	case klog.LevelDebug:
		zl.Debug("", data...)
	case klog.LevelInfo:
		zl.Info("", data...)
	case klog.LevelWarn:
		zl.Warn("", data...)
	case klog.LevelError:
		zl.Error("", data...)
	}
	return nil
}
