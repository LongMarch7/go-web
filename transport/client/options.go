package client

import (
    "context"
    "github.com/go-kit/kit/sd"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "time"
)

func EtcdServer(etcd string) COption {
    return func(o *ClientOpt) {
        o.etcdServer = etcd
    }
}

func Prefix(prefix string) COption {
    return func(o *ClientOpt) {
        o.prefix = prefix
    }
}

func Ctx(ctx context.Context) COption {
    return func(o *ClientOpt) {
        o.ctx = ctx
    }
}

func Mux(mux *runtime.ServeMux) COption {
    return func(o *ClientOpt) {
        o.mux = mux
    }
}

func DialTimeout(dialTimeout  time.Duration) COption {
    return func(o *ClientOpt) {
        o.dialTimeout = dialTimeout
    }
}

func DialKeepAlive(dialKeepAlive  time.Duration) COption {
    return func(o *ClientOpt) {
        o.dialKeepAlive = dialKeepAlive
    }
}

func Factory(factory  sd.Factory) COption {
    return func(o *ClientOpt) {
        o.factory = factory
    }
}

func RetryTime(retryTime  time.Duration) COption {
    return func(o *ClientOpt) {
        o.retryTime = retryTime
    }
}

func RetryCount(retryCount  int) COption {
    return func(o *ClientOpt) {
        o.retryCount = retryCount
    }
}

func RegisterGrpc(registerGrpc  RegisterHandlerClient) COption {
    return func(o *ClientOpt) {
        o.registerGrpc = registerGrpc
    }
}

func Extend(extend  interface{}) COption {
    return func(o *ClientOpt) {
        o.extend = extend
    }
}

func Manager(manager  interface{}) COption {
    return func(o *ClientOpt) {
        o.manager = manager
    }
}

func HystrixTimeout(hystrixTimeout  int) COption {
   return func(o *ClientOpt) {
       o.hystrixTimeout = hystrixTimeout
   }
}

func HystrixErrorPercentThreshold(hystrixErrorPercentThreshold  int) COption {
    return func(o *ClientOpt) {
        o.hystrixErrorPercentThreshold = hystrixErrorPercentThreshold
    }
}

func HystrixSleepWindow(hystrixSleepWindow  int) COption {
    return func(o *ClientOpt) {
        o.hystrixSleepWindow = hystrixSleepWindow
    }
}

func HystrixMaxConcurrentRequests(hystrixMaxConcurrentRequests  int) COption {
    return func(o *ClientOpt) {
        o.hystrixMaxConcurrentRequests = hystrixMaxConcurrentRequests
    }
}

func HystrixRequestVolumeThreshold(hystrixRequestVolumeThreshold  int) COption {
    return func(o *ClientOpt) {
        o.hystrixRequestVolumeThreshold = hystrixRequestVolumeThreshold
    }
}