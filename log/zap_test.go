package log

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestZap(t *testing.T) {
	logger := NewZapLogger(&Options{
		Level:       zapcore.InfoLevel,
		ServiceName: "",
		Skip:        2,
		Writer:      NewFileWriter(&FileOption{
			Filename: "%Y-%m-%d.log",
			MaxSize:  20,
			MaxAge:   0,
		}),
	})


	logHelper := log.NewHelper(logger)

	logHelper.Info("111", 123)
	logHelper.Debugw("debug", 123)
}
