package userAccess

import (
	"github.com/bighuangbee/gokit/tools/coper"
	"github.com/bighuangbee/gokit/userAccess/jwt"
	"time"
)

const JWT_Encrtpy = "7#23e!GJd&AAz=13Da%6"

type TokenJWT struct{
}

func NewTokenJWT() IToken {
	return &TokenJWT{}
}

func (this *TokenJWT) Generate(claims *UserClaims) (string, error) {
	data := make(map[string]interface{})
	coper.Copy(&data, claims)

	token, err := jwt.Generate(time.Duration(claims.JwtClaims.ExpiresAt)*time.Second, []byte(JWT_Encrtpy), data)
	return token.AccessToken, err
}

func (this *TokenJWT)Decode(tokenStr string)(claims *UserClaims, err error){
	token, err := jwt.Parse(tokenStr, []byte(JWT_Encrtpy))
	if err != nil{
		return nil, err
	}

	claims = &UserClaims{}
	coper.Copy(claims, token.Data)
	return claims, err
}
