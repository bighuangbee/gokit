package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	//Level default INFO
	Level zapcore.Level
	//ServiceName serviceName
	ServiceName string
	//调用者堆栈层级
	Skip int
	//日志存储器
	Writer io.Writer
}

type ZapLogger struct {
	*zap.SugaredLogger
	Sync func() error
}

func NewZapLogger(opt *Options) *ZapLogger {
	if opt == nil {
		panic("Options is required")
	}

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		beijingT := t.In(time.FixedZone("CST", 8*60*60))
		enc.AppendString(beijingT.Format("2006-01-02 15:04:05"))
	}

	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder, // Use custom time encoder
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapopts := []zap.Option{}
	if opt.Skip != 0 {
		encoder.CallerKey = "caller"
		zapopts = append(zapopts, zap.AddCallerSkip(opt.Skip))
	}

	// Log storage
	ws := []zapcore.WriteSyncer{
		zapcore.AddSync(os.Stdout),
	}
	if opt.Writer != nil {
		ws = append(ws, zapcore.AddSync(opt.Writer))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), // Encoder configuration
		zapcore.NewMultiWriteSyncer(ws...),
		zap.NewAtomicLevelAt(opt.Level),
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

func (l *ZapLogger) Logw(level zapcore.Level, msg string, keyvals ...interface{}) {
	if len(keyvals) == 0 {
		l.SugaredLogger.Desugar().Log(level, msg)
	} else if len(keyvals)%2 == 0 {
		var data []zap.Field
		for i := 0; i < len(keyvals); i += 2 {
			data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
		}
		l.SugaredLogger.Desugar().Log(level, msg, data...)
	} else {
		l.SugaredLogger.Desugar().Log(level, msg, zap.Any(fmt.Sprint(keyvals[0]), ""))
	}
}

func (l *ZapLogger) WithCtx(ctx context.Context) Logger {
	return l.WithCtx(ctx)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.Logw(zapcore.DebugLevel, "", args...)
}
func (l *ZapLogger) Info(args ...interface{}) {
	l.Logw(zapcore.InfoLevel, "", args...)
}
func (l *ZapLogger) Warn(args ...interface{}) {
	l.Logw(zapcore.WarnLevel, "", args...)
}
func (l *ZapLogger) Error(args ...interface{}) {
	l.Logw(zapcore.ErrorLevel, "", args...)
}

func (l *ZapLogger) Debugf(msg string, args ...interface{}) {
	l.Logw(zapcore.DebugLevel, msg, args...)
}
func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	fmt.Println("=======Infof", len(args))
	l.Logw(zapcore.InfoLevel, msg, args)
}
func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.Logw(zapcore.WarnLevel, msg, args...)
}
func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.Logw(zapcore.ErrorLevel, msg, args...)
}

func (l *ZapLogger) Debugw(msg string, args ...interface{}) {
	l.Logw(zapcore.DebugLevel, msg, args...)
}
func (l *ZapLogger) Infow(msg string, args ...interface{}) {
	l.Logw(zapcore.InfoLevel, msg, args...)
}
func (l *ZapLogger) Warnw(msg string, args ...interface{}) {
	l.Logw(zapcore.WarnLevel, msg, args...)
}
func (l *ZapLogger) Errorw(msg string, keyvals ...interface{}) {
	l.Logw(zapcore.ErrorLevel, msg, keyvals...)
}

// ==================文件存储日志==================
type FileOption struct {
	Filename string //指定日志存储位置, 如按日期分割, filename包含%Y-%m-%d.log，/opt/logs/xxx/%Y-%m-%d.log
	MaxSize  int    //日志的最大大小（M）
	MaxAge   int    //日志文件存储最大天数
}

func NewFileWriter(file *FileOption) io.Writer {
	if file.Filename == "" {
		return nil
	}

	o := []rotatelogs.Option{
		rotatelogs.WithRotationTime(24 * time.Hour), //分割时长,默认最少60s
		//rotatelogs.WithLinkName(file.Filename), //生成软链，指向最新日志文件
	}

	//保存天数
	if file.MaxAge > 0 {
		o = append(o, rotatelogs.WithMaxAge(time.Hour*24*time.Duration(file.MaxAge)))
	}
	//大于这个容量创建新的日志文件
	if file.MaxSize > 0 {
		o = append(o, rotatelogs.WithRotationSize(1024*1024*int64(file.MaxSize)))
	}

	hook, err := rotatelogs.New(file.Filename, o...)
	if err != nil {
		panic("NewFileWriter " + err.Error())
	}

	return hook
}
