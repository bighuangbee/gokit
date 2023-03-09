/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package userAccess

import (
	"errors"
	"fmt"
	"github.com/bighuangbee/gokit/tools/cache"
	"time"
)


type IUserAccess interface {
	//签发token
	Issue(*UserClaims) (token string, err error)
	//验证token
	Validate(token string)(*UserClaims, error)
	//注销token
	Logout(username string)(error)
}


func New(store cache.Cache, loginExpire time.Duration) *UserAccess {
	return &UserAccess{Store: store, AccessToken: NewAccessToken(), LoginExpire: loginExpire}
}

// implement IUserAccess interface
type UserAccess struct {
	Store       cache.Cache
	AccessToken IAccessToken
	LoginExpire time.Duration //登陆有效期
}

//签发token
func (userService *UserAccess)Issue(user *UserClaims)(string, error){
	claims := user
	claims.JwtClaims.ExpiresAt = time.Now().Add(time.Minute * userService.LoginExpire).Unix()

	token, err := userService.AccessToken.Generate(claims)
	if err != nil{
		return "", err
	}

	//err = userService.Store.SetEntity(userService.AccessToken.StoreKey(claims), claims, userService.LoginExpire)
	return token, err
}

//验证token
func (userService *UserAccess) Validate(token string) (*UserClaims, error){
	if token == ""{
		return nil, errors.New("token not allow empty.")
	}

	user, err := userService.AccessToken.Decode(token)
	if err != nil{
		return nil, err
	}

	var validUser UserClaims
	//err = userService.Store.GetEntity(userService.AccessToken.StoreKey(user), &validUser)

	fmt.Println("--validUser", validUser, user)
	return nil, errors.New("token invalid")
}

//注销token
func (userService *UserAccess)Logout(token string)error{
	user, err := userService.AccessToken.Decode(token)
	if err != nil{
		return err
	}

	fmt.Println("p---user", user)
	//return userService.Store.Del(userService.AccessToken.StoreKey(user))
	return nil
}
