package userAccess

import (
	"context"
	"errors"
	"github.com/bighuangbee/gokit/storage/cache"
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

func New(store cache.ICache, loginExpire time.Duration) *UserAccess {
	return &UserAccess{
		Store: store,
		Token: NewToken(),
		LoginExpire: loginExpire,
		Ctx: context.Background(),
	}
}

// implement IUserAccess interface
type UserAccess struct {
	//储存token
	Store       cache.ICache
	//生成token，解析token
	Token       IToken
	//登陆有效期
	LoginExpire time.Duration
	Ctx context.Context
}

//签发token
func (userService *UserAccess)Issue(user *UserClaims)(string, error){
	user.JwtClaims.ExpiresAt = time.Now().Add(time.Minute * userService.LoginExpire).Unix()
	token, err := userService.Token.Generate(user)
	if err != nil{
		return "", err
	}

	user.Token = token
	err = userService.Store.SetEntity(userService.Ctx, user.StoreKey(), user, userService.LoginExpire)
	return token, err
}

//验证token
func (userService *UserAccess) Validate(token string) (*UserClaims, error){
	if token == ""{
		return nil, errors.New("token not allow empty.")
	}

	user, err := userService.Token.Decode(token)
	if err != nil{
		return nil, err
	}

	var validUser UserClaims
	if err := userService.Store.GetEntity(userService.Ctx, user.StoreKey(), &validUser); err != nil{
		return nil, errors.New("token invalid")
	}

	if validUser.Token == token{
		return user, nil
	}

	return nil, errors.New("validate error.")
}

//注销token
func (userService *UserAccess)Logout(token string)error{
	user, err := userService.Token.Decode(token)
	if err != nil{
		return err
	}
	return userService.Store.Del(userService.Ctx, user.StoreKey())
}
