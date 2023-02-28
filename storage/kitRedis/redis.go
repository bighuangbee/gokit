package kitRedis

import (
	"context"
	"time"

	"github.com/bighuangbee/gokit/log"
	"github.com/go-redis/redis/v8"
)

const (
	// 单节点模式
	ModeSingle = "single"
	// 集群模式
	ModeCluster = "cluster"
)

// default settings
const (
	_PoolSize     = 50
	_MinIdleConns = 1
	_IdleTimeout  = time.Hour
)

type Client interface {
	// Mode 连接方式
	Mode() string
	// RClient 返回*redis.Client或*redis.ClusterClient
	RClient() interface{}

	// redis 方法，自行添加

	Close() error
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd

	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
	XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd
	XLen(ctx context.Context, stream string) *redis.IntCmd
	XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd
	XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd
	XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd
	XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *redis.XMessageSliceCmd
	XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd
	XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd
	XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd
	XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd
	XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd
	XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd
	XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd
	XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd
	XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd
	XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd
	XPending(ctx context.Context, stream, group string) *redis.XPendingCmd
	XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd
	XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd
	XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd

	XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd
	XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd
	XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd
	XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd
	XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd
	XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd
	XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd

	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZCount(ctx context.Context, key, min, max string) *redis.IntCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd
	ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd
	ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd

	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd
	HKeys(ctx context.Context, key string) *redis.StringSliceCmd
	HLen(ctx context.Context, key string) *redis.IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd
	HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd
	HVals(ctx context.Context, key string) *redis.StringSliceCmd
	HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd
}

// Options options
type Options struct {
	// Mode 连接方式.
	// 单节点:"single",集群:"cluster",默认single.
	Mode string
	// Addr e.g. "127.0.0.1:6739" 单节点模式地址
	Addr string
	// ClusterAddrs 集群模式地址
	ClusterAddrs []string
	// Password optional
	Password string
	// DB default 0
	DB int
	// Logger logger, 命令将以DEBUG的level打印
	// example:
	// Logger: hiZap.New(&hiZap.Options{
	// 	Level: zapcore.DebugLevel,
	// 	Skip:  4,
	// }),
	Logger log.Logger
}

func New(opt *Options) (Client, error) {
	c := &baseClient{}
	// 单节点
	if opt.Mode != ModeCluster {
		c.mode = ModeSingle
		return newSingleClient(c, opt)
	}
	// 集群
	c.mode = ModeCluster
	return newClusterClient(c, opt)
}

type baseClient struct {
	mode string
}

func (c *baseClient) Mode() string {
	return c.mode
}

// 单节点连接
type singleClient struct {
	*baseClient
	*redis.Client
}

// 单节点连接
func newSingleClient(base *baseClient, opt *Options) (*singleClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         opt.Addr,
		Password:     opt.Password,
		DB:           opt.DB,
		PoolSize:     _PoolSize,
		MinIdleConns: _MinIdleConns,
		IdleTimeout:  _IdleTimeout,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	//client.AddHook(redisotel.NewTracingHook())

	if opt.Logger != nil {
		client.AddHook(&hook{Logger: opt.Logger})
	}
	c := &singleClient{baseClient: base, Client: client}
	return c, nil
}

func (c *singleClient) RClient() interface{} {
	return c.Client
}

// 集群连接
type clusterClient struct {
	*baseClient
	*redis.ClusterClient
}

// 集群连接
func newClusterClient(base *baseClient, opt *Options) (*clusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        opt.ClusterAddrs,
		Password:     opt.Password,
		PoolSize:     _PoolSize,
		MinIdleConns: _MinIdleConns,
		IdleTimeout:  _IdleTimeout,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	//client.AddHook(redisotel.NewTracingHook())

	if opt.Logger != nil {
		client.AddHook(&hook{Logger: opt.Logger})
	}
	c := &clusterClient{baseClient: base, ClusterClient: client}
	return c, nil
}

func (c *clusterClient) RClient() interface{} {
	return c.ClusterClient
}

type hook struct {
	Logger log.Logger
}

func (h *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	h.Logger.WithCtx(ctx).Debugw("redisStart",
		"cmd", cmd.String(),
	)
	return ctx, nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	h.Logger.WithCtx(ctx).Debugw("redisEnd",
		"cmd", cmd.String(),
	)
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	cmdsStr := make([]string, len(cmds))
	for i, v := range cmds {
		cmdsStr[i] = v.String()
	}
	h.Logger.WithCtx(ctx).Debugw("redisStart",
		"cmd", cmdsStr,
	)
	return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	cmdsStr := make([]string, len(cmds))
	for i, v := range cmds {
		cmdsStr[i] = v.String()
	}
	h.Logger.WithCtx(ctx).Debugw("redisStart",
		"cmd", cmdsStr,
	)
	return nil
}
