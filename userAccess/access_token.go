/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package userAccess

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type IAccessToken interface {
	//生成token
	Generate(claims *UserClaims)(token string, err error)
	//解析token
	Decode(token string)(claims *UserClaims, err error)
}

func NewAccessToken() IAccessToken {
	//return NewAccessTokenNormal()
	return NewAccessTokenJWT()
}


type UserClaims struct {
	Account string   `json:"account"` //账号
	UserName string   `json:"user_name"` //用户名称
	CorpId      int32    `json:"corp_id"` //企业 ID
	JwtClaims jwt.StandardClaims `json:"jwt"`
}

func NewUserClaims(userName string, account string, CorpId int32, userId string, expire time.Duration) *UserClaims {
	return &UserClaims{
		UserName: userName,
		Account: account,
		CorpId: CorpId,
		JwtClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Id:        userId,
			IssuedAt: time.Now().Unix(),
			Issuer: "bighuangbee",
	}}
}
