package kitRedis

import (
	"context"
	"testing"

	"git.hiscene.net/hifoundry/go-kit/constance/hiCommon"
	"git.hiscene.net/hifoundry/go-kit/log/hiZap"
	"go.uber.org/zap/zapcore"
)

func Test_New(t *testing.T) {
	_, err := New(&Options{
		Addr: "192.168.23.68:30096",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_Redis(t *testing.T) {
	client, err := New(&Options{
		Addr: "192.168.23.68:30096",
		Logger: hiZap.New(&hiZap.Options{
			Level: zapcore.DebugLevel,
			Skip:  4,
		}),
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ctx := context.WithValue(context.Background(), hiCommon.TIDKey, "ttttt")
	t.Run("Normal", func(t *testing.T) {
		_, err := client.Get(ctx, "gokit_test").Result()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})
	t.Run("Error", func(t *testing.T) {
		_, err := client.Get(ctx, "wwwwwwwwww").Result()
		if err == nil {
			t.Fail()
		}
	})
}
