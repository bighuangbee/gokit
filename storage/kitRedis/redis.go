/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package kitRedis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var RedisOnce sync.Once
var client *redis.Client

func NewRedis(address string, auth string, dbIndex int, logger log.Logger)(r *redis.Client, err error){
	RedisOnce.Do(func() {
		logger.Log(log.LevelInfo, "NewRedis", "", "address", address, "dbIndex", dbIndex)
		client = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: auth,
			DB:       dbIndex,
		})
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = client.Ping(ctx).Result()
		if err != nil{
			logger.Log(log.LevelInfo, "NewRedis err", err.Error())
			panic("NewRedis: "+ err.Error())
		}
	})

	return client, err
}

/**
 * @Description:
 * @param address 哨兵集群地址
 * @param auth
 * @param dbIndex
 */
func NewRedisFailover(address []string, auth string, dbIndex int)(r *redis.Client, err error){
	RedisOnce.Do(func() {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName: 	"mymaster",
			SentinelAddrs:  address,
			Password: 		auth,
			ReadTimeout: 	200 * time.Millisecond,
			WriteTimeout: 	200 * time.Millisecond,
			DB: 			dbIndex,
		})
		ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
		_, err = client.Ping(ctx).Result()
	})
	return client, err
}
