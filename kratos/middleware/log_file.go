package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
	nethttp "net/http"
	"strings"
	"time"
)

// CheckTokenMiddleWare Check Token middleware
func LogFile(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {

				var (
					code      int32
					reason    string
					kind      string
					operation string
				)

				startTime := time.Now()
				kind = tr.Kind().String()
				operation = tr.Operation()

				httpTr := tr.(*http.Transport)

				reply, err = handler(ctx, req)
				if se := errors.FromError(err); se != nil {
					code = se.Code
					reason = se.Reason
				}

				body, _ := json.Marshal(req)

				logger.Log(log.LevelInfo, "server", "request_log",
					"component", kind,
					"url", httpTr.Request().URL.String(),
					"operation", operation,
					"method", httpTr.Request().Method,
					"body", string(body),
					"IP", getIP(httpTr.Request()),
					"code", code,
					"reason", reason,
					"latency", fmt.Sprintf("%dms", time.Since(startTime).Milliseconds()))
			}
			return
		}
	}
}

func getIP(r *nethttp.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return ""
	}
	return remoteIP.String()
}

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}
