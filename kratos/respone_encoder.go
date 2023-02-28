package kratos

import (
	"fmt"
	kjson "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	netHttp "net/http"
)

type ResponeError struct {
	Code      int32             `json:"code,omitempty"`
	DetailMsg string            `json:"detailMsg,omitempty"`
	Message   string            `json:"message,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	RetCode   int32             `json:"retCode"`
}

func ErrorEncoder() http.ServerOption {
	return http.ErrorEncoder(func(w netHttp.ResponseWriter, r *netHttp.Request, err error) {
		se := errors.FromError(err)

		//获取编码器
		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(&ResponeError{
			DetailMsg: se.Reason,
			Message:   se.Message,
			Metadata:  se.Metadata,
			RetCode:   se.Code,
		})
		if err != nil {
			w.WriteHeader(netHttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		code := netHttp.StatusOK
		if se.Code == netHttp.StatusUnauthorized {
			code = netHttp.StatusUnauthorized
		}
		w.WriteHeader(code)
		w.Write(body)
		return
	})
}

const (
	defaultReason  = "SUCCESS"
	defaultMessage = "操作成功"
)

func SuccessEncoder() http.ServerOption {
	return http.ResponseEncoder(func(w netHttp.ResponseWriter, r *netHttp.Request, v interface{}) error {
		msg := GetMessage(w)
		if msg == "" {
			msg = defaultMessage
		}
		codec, _ := http.CodecForRequest(r, "Accept")
		// 枚举使用数字
		kjson.MarshalOptions.UseEnumNumbers = true
		data, err := codec.Marshal(v)


		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write([]byte(fmt.Sprintf(`{
			"code": %d,
			"detailMsg": "%s",
			"message": "%s",
			"data": %s
			}`, 0, defaultReason, msg, data)))
		return nil
	})
}
