package server
import (
    "context"
    "time"
    "google.golang.org/grpc"
)

func EtcdServer(etcd string) SOption {
    return func(o *ServerOpt) {
        o.EtcdServer = etcd
    }
}

func Prefix(prefix string) SOption {
    return func(o *ServerOpt) {
        o.Prefix = prefix
    }
}

func ServerAddr(server string) SOption {
    return func(o *ServerOpt) {
        o.ServerAddr = server
    }
}

func Ctx(ctx context.Context) SOption {
    return func(o *ServerOpt) {
        o.Ctx = ctx
    }
}

func DialTimeout(dialTimeout  time.Duration) SOption {
    return func(o *ServerOpt) {
        o.DialTimeout = dialTimeout
    }
}

func DialKeepAlive(dialKeepAlive  time.Duration) SOption {
    return func(o *ServerOpt) {
        o.DialKeepAlive = dialKeepAlive
    }
}

func MaxThreadCount(maxThreadCount  string) SOption {
    return func(o *ServerOpt) {
        o.MaxThreadCount = maxThreadCount
    }
}

func NetType(netType  string) SOption {
    return func(o *ServerOpt) {
        o.NetType = netType
    }
}

func ServiceInit(init  func() interface {}) SOption {
    return func(o *ServerOpt) {
        o.ServiceStruct = init()
    }
}

type RegisterServer func(*grpc.Server, interface{})
func RegisterServiceFunc(register RegisterServer) SOption {
    return func(o *ServerOpt) {
        o.RegisterServerFunc = register
    }
}
