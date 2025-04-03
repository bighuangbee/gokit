package function

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 生成签名函数
func Sign(appSecret string, params map[string]interface{}) string {
	delete(params, "sign")

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var content strings.Builder
	for i, key := range keys {
		value := params[key]
		if value != nil {
			if i > 0 {
				content.WriteString("&")
			}
			content.WriteString(fmt.Sprintf("%s=%v", key, value))
		}
	}

	signString := content.String()
	finalSignString := strings.ToLower(signString) + "&" + appSecret

	return MD5(finalSignString)
}

func isTimestampValid(timestamp interface{}, maxDiff int64) bool {
	ts, ok := timestamp.(int64)
	if !ok {
		// 如果时间戳为字符串，则尝试转换为 int64
		if tsStr, ok := timestamp.(string); ok {
			tsInt, err := strconv.ParseInt(tsStr, 10, 64)
			if err != nil {
				return false
			}
			ts = tsInt
		} else {
			return false
		}
	}

	currentTime := time.Now().Unix()
	diff := currentTime - ts

	return diff >= -maxDiff && diff <= maxDiff
}

// 服务端签名验证函数
func VerifySign(appSecret string, params map[string]interface{}) error {
	if !isTimestampValid(params["timestamp"], 3) {
		return errors.New("时间戳校验失败")
	}

	clientSign, ok := params["sign"].(string)
	if !ok {
		return errors.New("缺少签名字段")
	}
	delete(params, "sign")

	// 对参数按 key 排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var content strings.Builder
	for i, key := range keys {
		value := params[key]
		if value != nil {
			if i > 0 {
				content.WriteString("&")
			}
			content.WriteString(fmt.Sprintf("%s=%v", key, value))
		}
	}

	// 生成签名
	signString := content.String()
	finalSignString := strings.ToLower(signString) + "&" + appSecret

	serverSign := MD5(finalSignString)
	if serverSign == clientSign {
		return nil
	}
	return errors.New("参数签名校验失败")
}
