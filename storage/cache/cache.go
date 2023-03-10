package cache

import (
	"context"
	"errors"
	"github.com/bighuangbee/gokit/storage/kitRedis"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type ICache interface {
	Set(ctx context.Context, key string, val interface{}, expire time.Duration)error
	Get(ctx context.Context, key string)interface{}
	SetEntity(ctx context.Context, key string, val interface{}, expire time.Duration)error
	GetEntity(ctx context.Context, key string, obj interface{})error
	Keys(ctx context.Context, key string)[]string
	Del(ctx context.Context, key ...string)error
	Incr(ctx context.Context, key string)error
	Expire(ctx context.Context, key string, expire time.Duration) error
	HSet(ctx context.Context, key string, val interface{}) error
	HGet(ctx context.Context, key string)interface{}
}



const EXPIRE_DEFAULT = time.Minute*5

const(
	CACHE_REDIS = iota
	CACHE_REDIS_CLUSTER
	CACHE_REDIS_FAILOVER
)

func New(t int, addr []string, auth string, dbIndex int, logger log.Logger) (engine ICache, err error) {
	if len(addr) == 0{
		return nil, errors.New("address empty.")
	}

	if logger == nil{
		logger = log.DefaultLogger
	}

	switch t {
		case CACHE_REDIS:
			client, err := kitRedis.NewRedis(addr[0], auth, dbIndex, logger)
			if err != nil{
				return nil, errors.New("CACHE_REDIS " + err.Error())
			}
			engine = NewCacheRedis(client)
		case CACHE_REDIS_FAILOVER:
			client, err := kitRedis.NewRedis(addr[0], auth, dbIndex, logger)
			if err != nil{
				return nil, errors.New("NewRedisFailover " + err.Error())
			}
			engine = NewCacheRedis(client)
		default:
			client, err := kitRedis.NewRedis(addr[0], auth, dbIndex, logger)
			if err != nil{
				return nil, errors.New("CACHE_REDIS " + err.Error())
			}
			engine = NewCacheRedis(client)
	}
	return engine, nil
}
