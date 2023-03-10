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

func NewCacheRedis(client *redis.Client) ICache {
	return &CacheRedis{Client: client}
}

func(this *CacheRedis)Set(ctx context.Context, key string, val interface{}, expire time.Duration)error{
	return this.Client.Set(ctx, key, val, expire).Err()
}

func(this *CacheRedis)Get(ctx context.Context, key string)interface{}{
	return this.Client.Get(ctx, key).Val()
}

func(this *CacheRedis)SetEntity(ctx context.Context, key string, val interface{}, expire time.Duration)error{
	data, err := json.Marshal(&val)
	if err != nil{
		return err
	}
	return this.Set(ctx, key, string(data), expire)
}

func(this *CacheRedis)GetEntity(ctx context.Context, key string, obj interface{})error{
	cmd := this.Client.Get(ctx, key)
	if cmd.Err() != nil{
		return errors.New("数据类型错误")
	}
	return json.Unmarshal([]byte(cmd.Val()), obj)
}

func(this *CacheRedis)Keys(ctx context.Context, key string)[]string{
	return this.Client.Keys(ctx, key).Val()
}

func(this *CacheRedis)Del(ctx context.Context, key ...string)error{
	return this.Client.Del(ctx, key...).Err()
}

func(this *CacheRedis)Incr(ctx context.Context, key string)error{
	return this.Client.Incr(ctx, key).Err()
}

func(this *CacheRedis)HGet(ctx context.Context, key string)interface{}{
	return this.Client.HGetAll(ctx, key).Val()
}

func(this *CacheRedis)HSet(ctx context.Context, key string, val interface{})error{
	return this.Client.HSet(ctx, key, val).Err()
}

func(this *CacheRedis)Expire(ctx context.Context, key string, expire time.Duration)error{
	return this.Client.Expire(ctx, key, expire).Err()
}

