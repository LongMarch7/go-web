package server
import (
    "context"
    "time"
    "google.golang.org/grpc"
)

func EtcdServer(etcd string) SOption {
    return func(o *ServerOpt) {
        o.etcdServer = etcd
    }
}

func Prefix(prefix string) SOption {
    return func(o *ServerOpt) {
        o.prefix = prefix
    }
}

func ServerAddr(server string) SOption {
    return func(o *ServerOpt) {
        o.serverAddr = server
    }
}

func Ctx(ctx context.Context) SOption {
    return func(o *ServerOpt) {
        o.ctx = ctx
    }
}

func DialTimeout(dialTimeout  time.Duration) SOption {
    return func(o *ServerOpt) {
        o.dialTimeout = dialTimeout
    }
}

func DialKeepAlive(dialKeepAlive  time.Duration) SOption {
    return func(o *ServerOpt) {
        o.dialKeepAlive = dialKeepAlive
    }
}

func MaxThreadCount(maxThreadCount  string) SOption {
    return func(o *ServerOpt) {
        o.maxThreadCount = maxThreadCount
    }
}

func NetType(netType  string) SOption {
    return func(o *ServerOpt) {
        o.netType = netType
    }
}

func ServiceInit(init  func() interface {}) SOption {
    return func(o *ServerOpt) {
        o.serviceStruct = init()
    }
}

type RegisterServer func(*grpc.Server, interface{})
func RegisterServiceFunc(register RegisterServer) SOption {
    return func(o *ServerOpt) {
        o.registerServerFunc = register
    }
}
