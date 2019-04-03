package server

import (
    "github.com/go-kit/kit/sd"
    "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/grpclog"
    "net"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "time"
    "context"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
    "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/go-kit/kit/log"
    zaplog "github.com/LongMarch7/go-web/plugin/zap-log"
    pluginzipkin "github.com/LongMarch7/go-web/plugin/zipkin"
)

var c chan os.Signal
var wg sync.WaitGroup

type IServer interface {
    init()
    Run()
}

type ServerOpt struct {
    consulAddr           string
    prefix               string
    serverAddr           string
    serverPort           int
    ctx                  context.Context
    dialTimeout          time.Duration
    dialKeepAlive        time.Duration
    netType              string
    maxThreadCount       string
    registerServerFunc   RegisterServer
    serviceStruct        interface{}
    advertiseAddress     string
    advertisePort        string
    logger               log.Logger
}

type Server struct {
    opts                  ServerOpt
    listenConnector       net.Listener
    registrar             sd.Registrar
}

type SOption func(o *ServerOpt)
func NewServer(opts ...SOption) IServer {
    return newServer(opts...)
}

func newServer(opts ...SOption) IServer {
    options := newOptions(opts...)
    s := &Server{
        opts: options,
    }
    s.init()
    return s
}
func newOptions(opts ...SOption) ServerOpt {
    opt := ServerOpt{
        consulAddr: "http://localhost:8500",
        prefix: "bookServer",
        serverAddr: "127.0.0.1",
        serverPort: 0,
        ctx: context.Background(),
        dialTimeout: time.Second * 3,
        dialKeepAlive: time.Second * 3,
        maxThreadCount: "1024",
        netType: "tcp",
        registerServerFunc: nil,
        serviceStruct: nil,
        advertiseAddress: "192.168.1.80",
        advertisePort: "10086",
        logger: zaplog.NewDefaultLogger(),
    }

    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func (s *Server)init(){


    ls, _ := net.Listen("tcp", s.opts.serverAddr+":"+strconv.Itoa(s.opts.serverPort))

    port := ls.Addr().(*net.TCPAddr).Port
    s.listenConnector = ls
    // 创建注册器
    config := RegisterConfig{
        consulAddress: s.opts.consulAddr,
        prefix: s.opts.prefix,
        service: s.opts.serverAddr,
        port: port,
        advertiseAddress: s.opts.advertiseAddress,
        advertisePort: s.opts.advertisePort,
        logger: s.opts.logger,
        maxThreadCount: s.opts.maxThreadCount,
    }
    registrar := Register(config)

    // 注册器启动注册
    s.registrar = registrar
}

func (s *Server)Run(){
    if s.opts.registerServerFunc == nil || s.opts.serviceStruct == nil {
        s.opts.logger.Log("RegisterServerFunc and ServiceStruct must set")
        return
    }
    c = make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    wg.Add(1)

    zip, _ := pluginzipkin.GetZipkinTracer(s.opts.prefix + "gateway_grpc_server")
    var opts []grpc.ServerOption
    if zip != nil {
        opts = append(opts,grpc_middleware.WithUnaryServerChain(
            otgrpc.OpenTracingServerInterceptor(zip, otgrpc.LogPayloads()),
        ),)
    }else{
        opts =[]grpc.ServerOption{grpc.UnaryInterceptor(grpc_transport.Interceptor)}
    }
    gs := grpc.NewServer(opts...)

    s.registrar.Register()
    defer func(){
        pluginzipkin.CloseZipkinCollector()
        s.registrar.Deregister()
        s.registrar = nil
        gs.Stop()
    }()
    go func() {
        endpoint := Endpoints{}
        health := DefaultHealthService{}
        endpoint.HealthEndpoint = MakeHealthEndpoint(&health)
        handler := MakeHttpHandler(s.opts.ctx, endpoint)
        http.ListenAndServe(":" + s.opts.advertisePort, handler)
    }()
    s.opts.registerServerFunc(gs, s.opts.serviceStruct)
    go gs.Serve(s.listenConnector)
    go Producer()
    wg.Wait()
}

func Producer(){
Loop:
    for{
        select {
        case s := <-c:
            grpclog.Error("Producer | get", s)
            break Loop
        default:
        }
        time.Sleep(500 * time.Millisecond)
    }
    wg.Done()
}