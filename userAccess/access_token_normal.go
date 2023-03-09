package userAccess

import (
	"errors"
	"fmt"
	"github.com/bighuangbee/gokit/tools"
	"github.com/bighuangbee/gokit/tools/crypto"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

const aeskey = "4Cd$1f!7#f13Aae1^eB30fDb3f72*3_N"

type AccessTokenNormal struct{
}

func NewAccessTokenNormal() IAccessToken {
	return &AccessTokenNormal{}
}

func (this *AccessTokenNormal) Generate(claims *UserClaims) (string, error) {
	str := fmt.Sprintf("%d_%s_%s_%s", claims.CorpId, claims.Account, claims.JwtClaims.Id, tools.MD5(aeskey+ "_basic-platform_" + time.Now().String()))
	return crypto.AesEncryptStr(str, aeskey)
}

func (this *AccessTokenNormal)Decode(token string)(claims *UserClaims, err error){
	str, err := crypto.AesDecryptStr(token, aeskey)
	if err != nil{
		return nil, err
	}

	decodeStr := strings.Split(str, "_")
	if len(decodeStr) > 3{
		corpId, _ := strconv.Atoi(decodeStr[0])
		account := decodeStr[1]
		userId := decodeStr[2]

		return &UserClaims{
			Account:       account,
			CorpId: int32(corpId),
			JwtClaims: jwt.StandardClaims{Id: userId},
		}, err
	}
	return nil, errors.New("Decode err")
}
