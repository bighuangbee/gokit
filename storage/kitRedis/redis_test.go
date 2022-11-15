package kitRedis

import (
	"context"
	"testing"
	"github.com/bighuangbee/gokit/constance"
	"github.com/bighuangbee/gokit/log/kitZap"
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
		Logger: kitZap.New(&kitZap.Options{
			Level: zapcore.DebugLevel,
			Skip:  4,
		}),
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ctx := context.WithValue(context.Background(), constance.TraceID, "ttttt")
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
