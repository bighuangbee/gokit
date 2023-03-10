/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package userAccess

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type IToken interface {
	//生成token
	Generate(claims *UserClaims)(token string, err error)
	//解析token
	Decode(token string)(claims *UserClaims, err error)
}

func NewToken() IToken {
	//return NewTokenNormal()
	return NewTokenJWT()
}


type UserClaims struct {
	From string   `json:"from"` //web wechat
	Account string   `json:"account"` //账号
	UserName string   `json:"user_name"` //用户名称
	CorpId int32    `json:"corp_id"` //企业 ID
	JwtClaims jwt.StandardClaims `json:"jwt"`
	Token string `json:"token"`
}

func (r *UserClaims) StoreKey() string {
	return fmt.Sprintf("%s_%d_%s_%s", r.From, r.CorpId, r.Account, r.JwtClaims.Id)
}

func NewUserClaims(userName string, account string, corpId int32, userId string, expire time.Duration) *UserClaims {
	return &UserClaims{
		UserName: userName,
		Account: account,
		CorpId: corpId,
		JwtClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Id:        userId,
			IssuedAt: time.Now().Unix(),
			Issuer: "bighuangbee",
	}}
}
