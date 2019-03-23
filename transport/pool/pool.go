package pool

import (
    "fmt"
    sq "github.com/LongMarch7/go-web/util/queue"
    "google.golang.org/grpc"
    "sync"
)

type Pool struct {
    Queue      *sq.EsQueue
    prefix     string
}

var poolManager = make( map[string] Pool, 128 )

type UpdatePool func(string, []string, uint32)

var lock sync.Mutex

func Update( prefix string, instances []string, count uint32){
    go func(){
        lock.Lock()
        for _, v := range instances {
            if pool, ok := poolManager[v]; ok{
                if pool.prefix != prefix{
                    for pool.Queue.Quantity() != 0 {
                        val, ok, _ := pool.Queue.Get()
                        if !ok {
                            val.(*grpc.ClientConn).Close()
                            val = nil
                        }
                    }
                    pool.prefix = prefix
                }
            }else{
                if count <= 64 {
                    count = 64
                }
                pool = Pool{
                    Queue: sq.NewQueue(count),
                    prefix: prefix,
                }
                poolManager[v] = pool
            }
            fmt.Println( "prefix = ", prefix,"instance = ",v)
        }
        lock.Unlock()
    }()
}

func GetConnect(key string) (Pool, bool){
    pool, ok := poolManager[key]
    return pool, ok
}

func Destroy(){
    for _,value :=range poolManager{
        for value.Queue.Quantity() != 0 {
           val, ok, _ := value.Queue.Get()
           if !ok {
               val.(*grpc.ClientConn).Close()
               val = nil
           }
        }
        value.Queue = nil
    }
    poolManager = nil
}