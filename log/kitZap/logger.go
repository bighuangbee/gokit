package kitZap

import (
	"context"
	"fmt"
	"os"

	"github.com/bighuangbee/gokit/constance"
	"github.com/bighuangbee/gokit/log"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/natefinch/lumberjack"
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
	//日志存储
	Storage *Storage
}

type Storage struct {
	Filename   string //指定日志存储位置
	MaxSize    int    //日志的最大大小（M）
	MaxBackups int    //日志的最大保存数量
	MaxAge     int    //日志文件存储最大天数
	Compress   bool   //是否压缩
}

func NewStorage() *Storage {
	return &Storage{Filename: "/opt/logs.log", MaxSize: 10, MaxBackups: 5, MaxAge: 30, Compress: false}
}

type ZapLogger struct {
	*zap.SugaredLogger
	Sync func() error
}

func New(opt *Options) *ZapLogger {
	if opt == nil {
		panic("Options is required")
	}

	encoder := zapcore.EncoderConfig{
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
	}

	zapopts := []zap.Option{}
	if opt.Skip != 0 {
		zapopts = append(zapopts, zap.AddCallerSkip(opt.Skip))
	}
	if opt.Storage == nil {
		opt.Storage = NewStorage()
	}

	level := zap.NewAtomicLevelAt(opt.Level)
	var core zapcore.Core
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(getLogWriter(opt.Storage))), // 打印到控制台和文件
		level,
	)

	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(opt.Skip),
		zap.Development())

	if opt.ServiceName != "" {
		zapLogger = zapLogger.With(zap.String("service", opt.ServiceName))
	}
	return &ZapLogger{
		SugaredLogger: zapLogger.Sugar(),
		Sync:          zapLogger.Sync,
	}
}

func (l *ZapLogger) WithCtx(ctx context.Context) log.Logger {
	r := &ZapLogger{}

	v := tracing.TraceID()
	tid := v(ctx).(string)
	if tid != "" {
		r.SugaredLogger = l.SugaredLogger.With(zap.String(constance.TraceID, v(ctx).(string)))
	} else {
		r.SugaredLogger = l.SugaredLogger
	}
	return r
}

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

// 日志自动切割
func getLogWriter(s *Storage) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   s.Filename,
		MaxSize:    s.MaxSize,
		MaxBackups: s.MaxBackups,
		MaxAge:     s.MaxAge,
		Compress:   s.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
