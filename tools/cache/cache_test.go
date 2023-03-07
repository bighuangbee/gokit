package cache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)


//var addr = []string{"localhost:6379", "localhost:6380", "localhost:6381"}
var addr = []string{"192.168.80.94:6379"}
var passwd = ""
var  index, _ = strconv.Atoi("0")

func TestCache(t *testing.T) {
	cache := New(CACHE_REDIS, addr, passwd, index)

	cache.Set("name", 10088, EXPIRE_DEFAULT*time.Second)
	fmt.Println("1 name:", cache.Get("name"))
	cache.Del("name", "name1")
	fmt.Println("2 name:", cache.Get("name"))
}

func TestCacheEntity(t *testing.T) {
	cache := New(CACHE_REDIS, addr, passwd, index)

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

	cache.SetEntity(cacheKey, user, EXPIRE_DEFAULT*time.Second)
	fmt.Println("1 :", cache.Get(cacheKey))

	var userStrust User
	cache.GetEntity(cacheKey, &userStrust)
	fmt.Println("2 :", userStrust)

	var userMap = make(map[string]interface{})
	cache.GetEntity(cacheKey, &userMap)
	fmt.Println("3 :", userMap)

}

func TestHSET(t *testing.T) {
	cache := New(CACHE_REDIS, addr, passwd, index)

	cache.Del("userinfo")

	userinfo := make(map[string]interface{})
	userinfo["name"] = "wiehua"
	userinfo["age"] = 11

	err := cache.HSet("userinfo", userinfo)
	fmt.Println("2 userinfo:", err,cache.HGet("userinfo"))

	err = cache.Expire("userinfo", 20*time.Second)
}

