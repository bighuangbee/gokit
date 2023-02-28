package kitRedis

import (
	"context"
	"github.com/bighuangbee/gokit/constance"
	"github.com/bighuangbee/gokit/log"
	"go.uber.org/zap/zapcore"
	"testing"
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
		Logger: log.NewZapLogger(&log.Options{
			Level:       zapcore.DebugLevel,
			ServiceName: "",
			Skip:        2,
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
