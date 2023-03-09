package jwtToken

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtToken struct {
	duration     time.Duration          // 有效时长 默认 1 分钟
	Data         map[string]interface{} // 数据
	AccessToken  string                 // 请求token
	RefreshToken string                 // 刷新token
	Key          []byte                 // 盐
}


// 生成token
func (t *jwtToken) generate() error {
	if t.AccessToken != "" {
		return nil
	}

	if t.Data == nil {
		t.Data = map[string]interface{}{}
	}
	if t.duration == 0 {
		t.duration = 2 * time.Hour
	}
	// 设置过期时间
	t.Data["exp"] = time.Now().Add(t.duration).Unix()

	if t.Key == nil {
		return errors.New("key无效为空")
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(t.Data))
	if token, err := at.SignedString(t.Key); err != nil {
		return err
	} else {
		t.AccessToken = token
	}

	return nil
}


// 解析
func (t *jwtToken) parse() error {
	if t.AccessToken == "" {
		return errors.New("Token无效为空")
	} else if t.Key == nil {
		return errors.New("Key无效为空")
	}

	claim, err := jwt.Parse(t.AccessToken, func(tt *jwt.Token) (interface{}, error) {
		return t.Key, nil
	})
	if err != nil {
		return err
	}

	t.Data = map[string]interface{}(claim.Claims.(jwt.MapClaims))
	return nil
}
