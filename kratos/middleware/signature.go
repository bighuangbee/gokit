package middleware

import (
	"context"
	"github.com/bighuangbee/gokit/function"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Signature(logger log.Logger, appSecret string, timeDiff int64) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if tr, ok := tr.(*http.Transport); ok {
					header := tr.Request().Header
					params := make(map[string]interface{})
					params["timestamp"] = header.Get("timestamp")
					params["nonce"] = header.Get("nonce")
					params["appId"] = header.Get("appId")
					params["sign"] = header.Get("sign")

					if err := function.VerifySign(appSecret, params, timeDiff); err != nil {
						logger.Log(log.LevelError, "err", err.Error(), "params", params)
						return nil, err
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
