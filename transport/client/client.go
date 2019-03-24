package client

import (
    "context"
    "fmt"
    "github.com/afex/hystrix-go/hystrix"
    "github.com/go-kit/kit/circuitbreaker"
    "github.com/go-kit/kit/sd"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/sd/lb"
    "google.golang.org/grpc"
    "time"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/transport/pool"
)

type BaseGatewayManager struct {
    Handler func(conn *grpc.ClientConn) error
    manager interface{}
}

func NewManager(manager interface{}) *BaseGatewayManager {
    return &BaseGatewayManager{
        manager: manager,
    }
}

type GrpcPoolManager struct{
    Opt       grpc.DialOption
    Extend    interface{}
}
type RegisterHandlerClient func( context.Context, *runtime.ServeMux, endpoint.Endpoint, interface{}) error
type IClient interface {
    init()
    DeRegister()
}

type ClientOpt struct {
    etcdServer      string
    prefix          string
    mux             *runtime.ServeMux
    ctx             context.Context
    dialTimeout     time.Duration
    dialKeepAlive   time.Duration
    registerGrpc    RegisterHandlerClient
    factory         sd.Factory
    retryTime       time.Duration
    retryCount      int
    extend          interface{}
    manager         interface{}
    hystrixTimeout                int
    hystrixMaxConcurrentRequests  int
    hystrixRequestVolumeThreshold int
    hystrixSleepWindow            int
    hystrixErrorPercentThreshold  int
}

type Client struct {
    opts      ClientOpt
    cancel    context.CancelFunc
    instancer *etcdv3.Instancer
    defaultEndpointer *sd.DefaultEndpointer
}


type COption func(o *ClientOpt)
func NewClient(opts ...COption) IClient {
    return newClient(opts...)
}

func newClient(opts ...COption) IClient {
    options := newOptions(opts...)
    c := &Client{
        opts: options,
    }
    c.init()
    return c
}
func newOptions(opts ...COption) ClientOpt {
    opt := ClientOpt{
        etcdServer: "127.0.0.1:2379",
        prefix: "/services/book/",
        ctx: context.Background(),
        mux: nil,
        dialTimeout: time.Second * 3,
        dialKeepAlive: time.Second * 3,
        factory: defaultReqFactory,
        retryTime: time.Second * 3,
        retryCount: 3,
        registerGrpc: nil,
        extend: nil,
        manager: nil,
        hystrixTimeout: 1000,
        hystrixErrorPercentThreshold: 50,
        hystrixSleepWindow: 5000,
        hystrixMaxConcurrentRequests: 100,
        hystrixRequestVolumeThreshold: 50,
    }

    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func (c *Client)init(){
    if c.opts.mux == nil || c.opts.registerGrpc == nil{
        fmt.Println("mux and grpc need set")
        return
    }
    commandName := c.opts.prefix + "hystrix"
    hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
        Timeout: c.opts.hystrixTimeout,
        ErrorPercentThreshold: c.opts.hystrixErrorPercentThreshold,
        SleepWindow: c.opts.hystrixSleepWindow,
        MaxConcurrentRequests: c.opts.hystrixMaxConcurrentRequests,
        RequestVolumeThreshold: c.opts.hystrixRequestVolumeThreshold,
    })
    breakerMw := circuitbreaker.Hystrix(commandName)

    options := etcdv3.ClientOptions{
        DialTimeout: c.opts.dialTimeout,
        DialKeepAlive: c.opts.dialKeepAlive,
    }
    //连接注册中心
    client, err := etcdv3.NewClient(c.opts.ctx, []string{c.opts.etcdServer}, options)
    if err != nil {
        panic(err)
    }
    logger := log.NewNopLogger()
    //创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
    instancer, err := etcdv3.NewInstancer(client, c.opts.prefix, logger, pool.Update)
    if err != nil {
        panic(err)
    }
    c.instancer = instancer

    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, c.opts.factory, logger)
    c.defaultEndpointer = endpointer

    //创建负载均衡器
    balancer := lb.NewRoundRobin(endpointer)

    /**
    我们可以通过负载均衡器直接获取请求的endPoint，发起请求
    reqEndPoint,_ := balancer.Endpoint()
    */

    /**
    也可以通过retry定义尝试次数进行请求
    */
    reqEndPoint := lb.Retry(c.opts.retryCount, c.opts.retryTime, balancer)

    reqEndPoint = breakerMw(reqEndPoint)

    ctx, cancel := context.WithCancel(c.opts.ctx)
    c.cancel = cancel

    if c.opts.manager == nil {
        c.opts.manager = &GrpcPoolManager{
            Opt: grpc.WithInsecure(),
            Extend: c.opts.extend,
        }
    }
    err = c.opts.registerGrpc(ctx, c.opts.mux, reqEndPoint, c.opts.manager)
    if err != nil {
        panic(err)
    }
}

func (c *Client)DeRegister(){
    c.cancel()
    c.defaultEndpointer.Close()
    c.defaultEndpointer = nil
    c.instancer.Stop()
    c.instancer = nil
    pool.Destroy()
}
