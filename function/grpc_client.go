package function

import (
	"log"
	"sync"
	"time"

	"github.com/shimingyah/pool"
	"google.golang.org/grpc"
)

var (
	// 初始化一次，key=grpc addr
	grpcOnceMap sync.Map
	// 复用连接，key=grpc addr
	poolConnMap sync.Map
)

// 支持复用连接、重连
// grpc服务地址或者 服务注册中心的地址
func GetGrpcClient(address string) *grpc.ClientConn {

	actual, isOld := grpcOnceMap.LoadOrStore(address, sync.Once{})
	if !isOld {
		grpcOnce := actual.(sync.Once)
		grpcOnce.Do(func() {
			poolConn, err := pool.New(address, pool.DefaultOptions)
			if err != nil {
				log.Fatalf("grpc failed to new pool: %v", err)
			}
			poolConnMap.Store(address, poolConn)
		})
	}
	var poolConnInMap interface{}
	var ok bool
	poolConnInMap, ok = poolConnMap.Load(address)
	for !ok {
		time.Sleep(time.Millisecond * 100)
		poolConnInMap, ok = poolConnMap.Load(address)
	}
	poolConn := poolConnInMap.(pool.Pool)
	conn, err := poolConn.Get()
	if err != nil {
		log.Fatalf("grpc failed to get conn: %v", err)
	}
	return conn.Value()
}
