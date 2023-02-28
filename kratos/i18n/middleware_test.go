package i18n

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type testTransport struct{ header headerCarrier }

func (tr *testTransport) Kind() transport.Kind            { return transport.KindHTTP }
func (tr *testTransport) Endpoint() string                { return "" }
func (tr *testTransport) Operation() string               { return "" }
func (tr *testTransport) RequestHeader() transport.Header { return tr.header }
func (tr *testTransport) ReplyHeader() transport.Header   { return tr.header }

func newTransportCtx(lang string) context.Context {
	hc := headerCarrier{}
	hc.Set("Accept-Language", lang)
	return transport.NewServerContext(context.Background(), &testTransport{hc})
}

func TestTranslator(t *testing.T) {
	bundle := New(Options{
		Paths: []string{
			"active.en.toml",
			"active.zh.toml",
		},
	})

	// languages 匹配的语言
	languages := []string{"chinese", "english"}

	// mockRespData 结构
	// msgId => [中文结果，英文结果]
	mockRespData := map[string][]string{
		"HelloWorld":    {"你好世界", "hello world"},
		"PersonInfo":    {"姓名:李四,年龄:16", "name:李四,age:16"},
		"PersonInfoErr": {"姓名:李四,年龄:<no value>", "name:李四,age:<no value>"},
	}

	tests := []struct {
		name  string
		err   error
		msgID string
		args  []interface{}
	}{
		{
			"Tr()：有效的msgId", nil, "HelloWorld", nil,
		},
		{
			"Tr()：无效的msgId", errors.New("not found in language"), "HelloPaul", nil,
		},
		{
			"Tr()：有效的msgId，实参=形参", nil, "PersonInfo", []interface{}{"李四", 16},
		},
		{
			"Tr()：有效的msgId，实参>形参", nil, "PersonInfo", []interface{}{"李四", 16, 18, 20},
		},
		{
			"Tr()：有效的msgId，实参<形参", nil, "PersonInfoErr", []interface{}{"李四"},
		},
	}

	for _, test := range tests {
		for idx, lang := range languages {
			name := fmt.Sprintf("%s，language：%s", test.name, lang)
			t.Run(name, func(t *testing.T) {
				handler := func(ctx context.Context, in interface{}) (interface{}, error) {
					tr, ok := FromContext(ctx)
					if !ok {
						return nil, errors.New("no md")
					}
					result, err := tr.Tr(test.msgID, test.args...)
					if err != nil {
						return nil, err
					}
					return result, nil
				}
				result, err := Translator(bundle)(handler)(newTransportCtx(lang), "foo")
				if err == nil {
					if result.(string) != mockRespData[test.msgID][idx] {
						t.Errorf("期望结果：%s，返回结果：%s", mockRespData[test.msgID][idx], result)
					}
					return
				}
				// panic 情况
				// 已知 msgId 无效会出现panic
				var errMsg string
				switch e := err.(type) {
				case *i18n.MessageNotFoundErr:
					errMsg = e.Error()
				case error:
					errMsg = e.Error()
				default:
					t.Log(e)
					errMsg = "未定义错误"
				}
				if test.err == nil {
					t.Errorf("期望error：%v，返回error：%s", test.err, errMsg)
				} else {
					if -1 == strings.Index(errMsg, test.err.Error()) {
						t.Errorf("期望error：%v，返回error：%s", test.err, errMsg)
					}
				}
			})
		}
	}
}
