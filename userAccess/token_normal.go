package userAccess

import (
	"errors"
	"fmt"
	"github.com/bighuangbee/gokit/function"
	"github.com/bighuangbee/gokit/function/crypto"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

const aeskey = "4Cd$1f!7#f13Aae1^eB30fDb3f72*3_N"

type TokenNormal struct {
}

func NewTokenNormal() IToken {
	return &TokenNormal{}
}

func (this *TokenNormal) Generate(claims *UserClaims) (string, error) {
	str := fmt.Sprintf("%d_%s_%s_%s_%d_%d", claims.CorpId, claims.Account, claims.JwtClaims.Id, function.MD5(aeskey+"_basic-platform_"+time.Now().String()), claims.JwtClaims.IssuedAt, claims.JwtClaims.ExpiresAt)
	return crypto.AesEncryptStr(str, aeskey)
}

func (this *TokenNormal) Decode(token string) (claims *UserClaims, err error) {
	str, err := crypto.AesDecryptStr(token, aeskey)
	if err != nil {
		return nil, err
	}

	decodeStr := strings.Split(str, "_")
	if len(decodeStr) > 3 {
		corpId, _ := strconv.Atoi(decodeStr[0])
		account := decodeStr[1]
		userId := decodeStr[2]

		return &UserClaims{
			Account:   account,
			CorpId:    int32(corpId),
			JwtClaims: jwt.StandardClaims{Id: userId},
		}, err
	}
	return nil, errors.New("Decode err")
}
