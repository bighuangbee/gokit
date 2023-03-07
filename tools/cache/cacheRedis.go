package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)


type CacheRedis struct {
	Client *redis.Client

}

var ctx = context.Background()

func NewCacheRedis(client *redis.Client) Cache {

	return &CacheRedis{Client: client}
}

func(this *CacheRedis)Set(key string, val interface{}, expire time.Duration)error{
	return Client.Set(ctx, key, val, expire).Err()
}

func(this *CacheRedis)Get(key string)interface{}{
	return Client.Get(ctx, key).Val()
}

func(this *CacheRedis)SetEntity(key string, val interface{}, expire time.Duration)error{
	data, err := json.Marshal(&val)
	if err != nil{
		return err
	}
	return this.Set(key, string(data), expire)
}

func(this *CacheRedis)GetEntity(key string, obj interface{})error{
	cmd := Client.Get(ctx, key)
	if cmd.Err() != nil{
		return errors.New("数据类型错误")
	}
	return json.Unmarshal([]byte(cmd.Val()), obj)
}

func(this *CacheRedis)Keys(key string)[]string{
	return Client.Keys(ctx, key).Val()
}

func(this *CacheRedis)Del(key ...string)error{
	return Client.Del(ctx, key...).Err()
}

func(this *CacheRedis)Incr(key string)error{
	return Client.Incr(ctx, key).Err()
}

func(this *CacheRedis)HGet(key string)interface{}{
	return Client.HGetAll(ctx, key).Val()
}

func(this *CacheRedis)HSet(key string, val interface{})error{
	return Client.HSet(ctx, key, val).Err()
}

func(this *CacheRedis)Expire(key string, expire time.Duration)error{
	return Client.Expire(ctx, key, expire).Err()
}

