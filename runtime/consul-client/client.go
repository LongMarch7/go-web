package client

import (
    "context"
    "github.com/afex/hystrix-go/hystrix"
    "github.com/go-kit/kit/circuitbreaker"
    "github.com/go-kit/kit/sd"
    "github.com/go-kit/kit/sd/lb"
    "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
    "github.com/hashicorp/consul/api"
    "google.golang.org/grpc"
    "time"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/runtime/pool"
    "github.com/LongMarch7/go-web/util/sd/consul"
    zaplog "github.com/LongMarch7/go-web/plugin/zap-log"
    pluginzipkin "github.com/LongMarch7/go-web/plugin/zipkin"
)

type GrpcPoolManager struct{
    Opt       []grpc.DialOption
    Extend    interface{}
}
type RegisterHandlerClient func( context.Context, *runtime.ServeMux, endpoint.Endpoint, interface{}) error
type IClient interface {
    init() bool
    DeRegister()
}

type ClientOpt struct {
    consulAddr      string
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
    passingOnly                   bool
    logger                        log.Logger
}

type Client struct {
    opts      ClientOpt
    cancel    context.CancelFunc
    defaultEndpointer *sd.DefaultEndpointer
}

var breakerMw endpoint.Middleware
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
        consulAddr: "http://localhost:8500",
        prefix: "bookServer",
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
        passingOnly: true,
        logger: zaplog.NewDefaultLogger(),
    }

    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func (c *Client)init() bool{
    if c.opts.mux == nil || c.opts.registerGrpc == nil{
        c.opts.logger.Log("mux and grpc need set")
        return false
    }
    commandName := c.opts.prefix + "hystrix"
    hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
        Timeout: c.opts.hystrixTimeout,
        ErrorPercentThreshold: c.opts.hystrixErrorPercentThreshold,
        SleepWindow: c.opts.hystrixSleepWindow,
        MaxConcurrentRequests: c.opts.hystrixMaxConcurrentRequests,
        RequestVolumeThreshold: c.opts.hystrixRequestVolumeThreshold,
    })
    breakerMw = circuitbreaker.Hystrix(commandName)

    var client consul.Client
    {
        consulConfig := api.DefaultConfig()

        consulConfig.Address = c.opts.consulAddr
        consulClient, err := api.NewClient(consulConfig)
        if err != nil {
            c.opts.logger.Log("api.NewClient error")
            return false
        }
        client = consul.NewClient(consulClient)
    }

    //创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
    consulTag := []string{"MicroServer", c.opts.prefix}
    pool.Init()
    instancer := consul.NewInstancer(client, c.opts.logger, c.opts.prefix, consulTag, c.opts.passingOnly, pool.Update)//pool.Update

    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, c.opts.factory,  c.opts.logger)
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

    ctx, cancel := context.WithCancel(c.opts.ctx)
    c.cancel = cancel

    zip , _ := pluginzipkin.GetZipkinTracer("gateway")
    dialOpts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
    if zip != nil {
        dialOpts = append(dialOpts,grpc.WithUnaryInterceptor(
            otgrpc.OpenTracingClientInterceptor(zip, otgrpc.LogPayloads()),
        ))
    }
    if c.opts.manager == nil {
        c.opts.manager = &GrpcPoolManager{
            Opt: dialOpts,
            Extend: c.opts.extend,
        }
    }
    err := c.opts.registerGrpc(ctx, c.opts.mux, reqEndPoint, c.opts.manager)
    if err != nil {
        c.opts.logger.Log("registerGrpc faild")
        return false
    }
    return true
}

func (c *Client)DeRegister(){
    c.cancel()
    c.defaultEndpointer.Close()
    c.defaultEndpointer = nil
    pluginzipkin.CloseZipkinCollector()
    pool.Destroy()
}
