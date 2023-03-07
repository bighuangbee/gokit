package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Set(key string, val interface{}, expire time.Duration)error
	Get(key string)interface{}
	SetEntity(key string, val interface{}, expire time.Duration)error
	GetEntity(key string, obj interface{})error
	Keys(key string)[]string
	Del(key ...string)error
	Incr(key string)error
	Expire(key string, expire time.Duration) error
	HSet(key string, val interface{}) error
	HGet(key string)interface{}
}



const EXPIRE_DEFAULT = 60 * 5 //过期时间 s

const(
	CACHE_REDIS = iota
	CACHE_REDIS_CLUSTER
	CACHE_REDIS_FAILOVER
)

func New(t int, addr []string, auth string, dbIndex int) Cache {
	if len(addr) == 0{
		return nil
	}

	var engine Cache
	switch t {
	case CACHE_REDIS:
		client, err := NewRedis(addr[0], auth, dbIndex)
		if err != nil{
			fmt.Println("CACHE_REDIS:",err)
			return nil
		}
		engine = NewCacheRedis(client)
	case CACHE_REDIS_FAILOVER:
		client, err := NewRedisFailover(addr, auth, dbIndex)
		if err != nil{
			fmt.Println("CACHE_REDIS_FAILOVER:",err)
			return nil
		}
		engine = NewCacheRedis(client)
	default:
		client, _ := NewRedis(addr[0], auth, dbIndex)
		engine = NewCacheRedis(client)
	}
	return engine
}
