package client

import (
    "context"
    "github.com/go-kit/kit/sd"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "time"
)

func EtcdServer(etcd string) COption {
    return func(o *ClientOpt) {
        o.EtcdServer = etcd
    }
}

func Prefix(prefix string) COption {
    return func(o *ClientOpt) {
        o.Prefix = prefix
    }
}

func Ctx(ctx context.Context) COption {
    return func(o *ClientOpt) {
        o.Ctx = ctx
    }
}

func Mux(mux *runtime.ServeMux) COption {
    return func(o *ClientOpt) {
        o.Mux = mux
    }
}

func DialTimeout(dialTimeout  time.Duration) COption {
    return func(o *ClientOpt) {
        o.DialTimeout = dialTimeout
    }
}

func DialKeepAlive(dialKeepAlive  time.Duration) COption {
    return func(o *ClientOpt) {
        o.DialKeepAlive = dialKeepAlive
    }
}

func Factory(factory  sd.Factory) COption {
    return func(o *ClientOpt) {
        o.Factory = factory
    }
}

func RetryTime(retryTime  time.Duration) COption {
    return func(o *ClientOpt) {
        o.RetryTime = retryTime
    }
}

func RetryCount(retryCount  int) COption {
    return func(o *ClientOpt) {
        o.RetryCount = retryCount
    }
}

func RegisterGrpc(registerGrpc  RegisterHandlerClient) COption {
    return func(o *ClientOpt) {
        o.RegisterGrpc = registerGrpc
    }
}

func Extend(extend  interface{}) COption {
    return func(o *ClientOpt) {
        o.Extend = extend
    }
}

func Manager(manager  interface{}) COption {
    return func(o *ClientOpt) {
        o.Manager = manager
    }
}

func HystrixTimeout(hystrixTimeout  int) COption {
   return func(o *ClientOpt) {
       o.HystrixTimeout = hystrixTimeout
   }
}

func HystrixErrorPercentThreshold(hystrixErrorPercentThreshold  int) COption {
    return func(o *ClientOpt) {
        o.HystrixErrorPercentThreshold = hystrixErrorPercentThreshold
    }
}

func HystrixSleepWindow(hystrixSleepWindow  int) COption {
    return func(o *ClientOpt) {
        o.HystrixSleepWindow = hystrixSleepWindow
    }
}

func HystrixMaxConcurrentRequests(hystrixMaxConcurrentRequests  int) COption {
    return func(o *ClientOpt) {
        o.HystrixMaxConcurrentRequests = hystrixMaxConcurrentRequests
    }
}

func HystrixRequestVolumeThreshold(hystrixRequestVolumeThreshold  int) COption {
    return func(o *ClientOpt) {
        o.HystrixRequestVolumeThreshold = hystrixRequestVolumeThreshold
    }
}