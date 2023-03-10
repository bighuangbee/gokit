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

var LoginExpire = time.Hour * 24 * 15

var addr = []string{"localhost:6379"}
var passwd = "A123!@#"
var index, _ = strconv.Atoi("0")


func TestUserAccess(t *testing.T) {
	logger := kratosLog.DefaultLogger
	c, _ := cache.New(cache.CACHE_REDIS, addr, passwd, index, logger)

	UserAccess	:= New(c, LoginExpire)

	user, _ 	:= mockValidateUser("10088", "123456")

	token, err := UserAccess.Issue(user)
	if err != nil {
		fmt.Println("UserAccess.Issue:",err)
		return
	}
	fmt.Println("token:",token)

	userValidate, err := UserAccess.Validate(token)
	if err != nil{
		fmt.Println("UserAccess.Validate err:", err)
	}else{
		fmt.Println("UserAccess.Validate succ, result :", userValidate)
	}

	//用错误的token进行验证
	userValidate2, err2 := UserAccess.Validate("1m6q+6Si/f4yC9nDaySnpdCoFnm/rZDgnTZpV43e00Pr6/I88QozMOT3fi0ZRlUX")
	if err2 != nil{
		fmt.Println("UserAccess.Validate 错误token, err:", err2)
	}else{
		fmt.Println("UserAccess.Validate 错误token, succ, result :", userValidate2)
	}

	//err = UserAccess.Logout(token)
	//fmt.Println("UserAccess.Logout", err)
	//
	//userValidate, err = UserAccess.Validate(token)
	//if err != nil{
	//	fmt.Println("UserAccess.Validate after logout err:", err)
	//}else{
	//	fmt.Println("UserAccess.Validate after logout succ, result :", userValidate)
	//}

}


func mockValidateUser(username string, password string) (*UserClaims, error){

	// 验证用户 ...

	return &UserClaims{
		From: "web",
		Account:   username,
		UserName:  "大黄蜂",
		CorpId:    1,
		JwtClaims: jwt.StandardClaims{Id: fmt.Sprintf("%d", 10086)},
	}, nil
}
