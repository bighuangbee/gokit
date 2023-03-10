package userAccess

import (
	"fmt"
	"testing"
	"time"
)

func Test_TokenJwt(t *testing.T) {
	token := NewToken()

	userClaims := NewUserClaims("weihua", "大黄蜂", 123, "456", time.Minute*10)
	tokenStr, err := token.Generate(userClaims)
	fmt.Println("token Generate", err, tokenStr)

	userClaimsDecode, err := token.Decode(tokenStr)
	if err != err {
		fmt.Println("token Decode err", err)
		return
	}
	fmt.Println("token Decode success:", userClaimsDecode, userClaimsDecode.JwtClaims.IssuedAt)
}


func Test_TokenNormal(t *testing.T) {
	token := NewTokenNormal()

	userClaims := NewUserClaims("大黄蜂", "weihua", 123, "456", time.Minute*10)
	tokenStr, err := token.Generate(userClaims)
	fmt.Println("token Generate", err, tokenStr)

	userClaimsDecode, err := token.Decode(tokenStr)
	if err != err {
		fmt.Println("token Decode err", err)
		return
	}
	fmt.Println("token Decode success:", *userClaimsDecode)
}
