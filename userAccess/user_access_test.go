/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package userAccess

import (
	"fmt"
	"github.com/bighuangbee/gokit/storage/cache"
	"github.com/dgrijalva/jwt-go"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	"strconv"
	"testing"
	"time"
)

var LoginExpire = 1440 * 2 * time.Minute //单位：分钟

var addr = []string{"localhost:26379"}
var passwd = "hiDronedb2020."
var  index, _ = strconv.Atoi("0")


func TestUserAccess(t *testing.T) {
	logger := kratosLog.DefaultLogger
	c, _ := cache.New(cache.CACHE_REDIS, addr, passwd, index, logger)

	//tokenStore 	:= NewTokenStore(nil, th)
	UserAccess 	:= New(c, LoginExpire)

	user, _ 	:= validateUser("10088", "123456")
	token, err := UserAccess.Issue(user)
	if err != nil {
		fmt.Println("UserAccess.Issue:",err)
		return
	}

	fmt.Println("new token:",token)

	userValidate, err := UserAccess.Validate(token)
	if err != nil{
		fmt.Println("userCliamsDecode1 err:", err)
	}else{
		fmt.Println("userCliamsDecode1 succ, result :", userValidate)
	}




	userValidate2, err2 := UserAccess.Validate("1m6q+6Si/f4yC9nDaySnpdCoFnm/rZDgnTZpV43e00Pr6/I88QozMOT3fi0ZRlUX")
	if err2 != nil{
		fmt.Println("userCliamsDecode2 err:", err2)
	}else{
		fmt.Println("userValidate2 succ, result :", userValidate2)
	}

}


func validateUser(username string, password string) (*UserClaims, error){

	// 验证用户 ...

	return &UserClaims{
		Account:   username,
		UserName:  "大黄蜂",
		CorpId:    1,
		JwtClaims: jwt.StandardClaims{Id: fmt.Sprintf("%d", 10086)},
	}, nil
}
