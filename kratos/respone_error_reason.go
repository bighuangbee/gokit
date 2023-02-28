package kratos

import (
	"context"
	"github.com/bighuangbee/gokit/kratos/i18n"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
)

var headerKey = "_response-msg_"

type PbErrorReason interface {
	String() string
}

// 返回 500 错误
func ReturnErr(ctx context.Context, reason PbErrorReason, args ...interface{}) *errors.Error {
	msgId := reason.String()
	return errors.New(http.StatusInternalServerError, msgId, i18n.Tr(ctx, msgId, args...))
}

func ReturnErr1(ctx context.Context, reason PbErrorReason, args ...interface{}) *errors.Error {
	msgId := reason.String()
	return errors.New(http.StatusInternalServerError, msgId, i18n.Tr(ctx, msgId, args...))
}

// ResponseErr BFF层 响应http response
func ResponseErr(ctx context.Context, errFn func(string, ...interface{}) *errors.Error, args ...interface{}) *errors.Error {
	err := errFn("")
	return errors.New(int(err.Code), err.Reason, i18n.Tr(ctx, err.Reason, args...)).WithMetadata(err.Metadata)
}

// ResponseErrWithError BFF层 响应http response
func ResponseErrWithError(ctx context.Context, e *errors.Error, args ...interface{}) *errors.Error {
	return errors.New(int(e.Code), e.Reason, i18n.Tr(ctx, e.Reason, args...))
}

// SetMessage 想要自定义返回message， 因框架问题，通过header船体
func SetMessage(ctx context.Context, msgId string, args ...interface{}) bool {
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		tr.ReplyHeader().Set(headerKey, i18n.Tr(ctx, msgId, args...))
	}
	return ok
}

// GetMessage 重写http.ResponseEncoder时，获取SetMessage设置的值
func GetMessage(w http.ResponseWriter) string {
	msg := w.Header().Get(headerKey)
	if msg != "" {
		w.Header().Del(headerKey)
	}
	return msg
}
