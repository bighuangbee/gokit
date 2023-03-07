/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var Client *redis.Client
var RedisOnce sync.Once

func NewRedis(address string, auth string, dbIndex int)(r *redis.Client, err error){
	RedisOnce.Do(func() {
		fmt.Println("NewRedis:",address, auth, dbIndex)
		Client = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: auth,
			DB:       dbIndex,
		})
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = Client.Ping(ctx).Result()
		if err != nil{
			panic("NewRedis: "+ err.Error())
		}
	})
	return Client, err
}

/**
 * @Description:
 * @param address 哨兵集群地址
 * @param auth
 * @param dbIndex
 */
func NewRedisFailover(address []string, auth string, dbIndex int)(r *redis.Client, err error){
	RedisOnce.Do(func() {
		Client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName: 	"mymaster",
			SentinelAddrs:  address,
			Password: 		auth,
			ReadTimeout: 	200 * time.Millisecond,
			WriteTimeout: 	200 * time.Millisecond,
			DB: 			dbIndex,
		})
		ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
		_, err = Client.Ping(ctx).Result()
	})
	return Client, err
}
