package kitZap

import (
	"context"
	"testing"

	"git.hiscene.net/hifoundry/go-kit/constance/hiCommon"
	"git.hiscene.net/hifoundry/go-kit/log"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
)

func console(logger log.Logger) {
	logger.Info("interface default")
	ctx := context.WithValue(context.Background(), hiCommon.TIDKey, "ttttt")
	logger.WithCtx(ctx).Debug("interface with ctx tid")
	logger.WithCtx(ctx).Errorf("interface errorf %s", "www")
}

func Test_Zap(t *testing.T) {
	logger := New(&Options{
		Level:       zapcore.DebugLevel,
		ServiceName: "srvvv",
	})
	t.Run("ZapInstance", func(t *testing.T) {
		logger.Info("default start")
		ctx := context.WithValue(context.Background(), hiCommon.TIDKey, "ttttt")
		logger.WithCtx(ctx).Debug("with ctx tid")
		logger.WithCtx(context.Background()).Debug("with ctx no tid")
		logger.Info("default end")
		logger.Infow("default json", "k1", "v1", "k2", "v2")
		logger.WithCtx(ctx).Infow("default json with ctx", "k1", "v1", "k2", "v2")
	})
	t.Run("Interface", func(t *testing.T) {
		console(logger)
	})
}

func Test_Kratos_Log(t *testing.T) {
	logger := New(&Options{
		Level:       zapcore.DebugLevel,
		ServiceName: "srvvv",
	})
	_ = logger.Log(klog.LevelInfo, "msg", "hhh")
}
