package jwtToken

import (
	"time"
)

// 通过生成token创建jwt数据
// d 有效时长   key 盐  data 数据
func Generate(d time.Duration, key []byte, data map[string]interface{}) (*jwtToken, error) {
	t := &jwtToken{
		duration: d,
		Key:      key,
		Data:     data,
	}
	if err := t.generate(); err != nil {
		return nil, err
	}
	return t, nil
}


// 从token解析jwt数据结构
func Parse(token string, key []byte) (*jwtToken, error) {
	t := &jwtToken{AccessToken: token, Key: key}
	if err := t.parse(); err != nil {
		return nil, err
	}
	return t, nil
}

