package kitRedis

import (
	"context"
	"time"
)

type MyRedis struct {
	Rdb    Client
	Prefix string
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisAdd(ctx context.Context, key, value string) error {
	err := t.Rdb.Set(ctx, t.Prefix+key, value, 0).Err()
	return err
}

// 手动添加过期时间的值
func (t *MyRedis) RedisAddAndExp(ctx context.Context, key, value string, exp time.Duration) error {
	err := t.Rdb.Set(ctx, t.Prefix+key, value, exp).Err()
	return err
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisGet(ctx context.Context, key string) (value string, err error) {
	return t.Rdb.Get(ctx, t.Prefix+key).Result()
}

// exist>0存在
func (t *MyRedis) RedisExist(ctx context.Context, key string) (exist int64, err error) {
	return t.Rdb.Exists(ctx, t.Prefix+key).Result()
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisHAdd(ctx context.Context, key, field, value string) error {
	err := t.Rdb.HSet(ctx, t.Prefix+key, field, value).Err()
	return err
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisHGet(ctx context.Context, key, field string) (value string, err error) {
	return t.Rdb.HGet(ctx, t.Prefix+key, field).Result()
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisHExist(ctx context.Context, key, field string) (exist bool, err error) {
	return t.Rdb.HExists(ctx, t.Prefix+key, field).Result()
}

// key不需要加 prefix，自动加
func (t *MyRedis) RedisIncr(ctx context.Context, key string) (value int64, err error) {
	return t.Rdb.Incr(ctx, t.Prefix+key).Result()
}

// 删除key
//
//	key不需要加 prefix，自动加
func (t *MyRedis) RedisDel(ctx context.Context, key string) (value int64, err error) {
	return t.Rdb.Del(ctx, t.Prefix+key).Result()
}
