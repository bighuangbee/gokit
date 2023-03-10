package cache

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)


//var addr = []string{"localhost:6379", "localhost:6380", "localhost:6381"}
var addr = []string{"localhost:6379"}
var passwd = "A123!@#"
var  index, _ = strconv.Atoi("0")

var ctx = context.Background()

func Test_Cache(t *testing.T) {
	cache, err := New(CACHE_REDIS, addr, passwd, index, nil)
	if err != nil{
		fmt.Println("err", err)
		return
	}

	cache.Set(ctx,"name", 10088, EXPIRE_DEFAULT)
	fmt.Println("cache.Get, name:", cache.Get(ctx,"name"))

	cache.Del(ctx,"name", "name1")
	fmt.Println("cache.Del, name:", cache.Get(ctx,"name"))
}

func TestCacheEntity(t *testing.T) {
	cache, err := New(CACHE_REDIS, addr, passwd, index, nil)
	if err != nil{
		fmt.Println("err", err)
		return
	}

	type User struct {
		Name string `json:"name"`
		Age int 	`json:"age"`
		Phone string
	}

	user := User{
		Name: "www",
		Age:  10,
		Phone: "13711112222",
	}
	cacheKey := "userObj"

	cache.SetEntity(ctx, cacheKey, &user, EXPIRE_DEFAULT)
	fmt.Println("SetEntity:", cache.Get(ctx,cacheKey))

	var userStrust User
	cache.GetEntity(ctx, cacheKey, &userStrust)
	fmt.Println("GetEntity:", userStrust)

	var userMap = make(map[string]interface{})
	cache.GetEntity(ctx, cacheKey, &userMap)
	fmt.Println("GetEntity map:", userMap)

}

func TestHSET(t *testing.T) {
	cache, err := New(CACHE_REDIS, addr, passwd, index, nil)
	if err != nil{
		fmt.Println("err", err)
		return
	}

	cache.Del(ctx, "userinfo")

	userinfo := make(map[string]interface{})
	userinfo["name"] = "wiehua"
	userinfo["age"] = 11

	err = cache.HSet(ctx,"userinfo", userinfo)
	fmt.Println("2 userinfo:", err,cache.HGet(ctx,"userinfo"))

	err = cache.Expire(ctx,"userinfo", 20*time.Second)
}

