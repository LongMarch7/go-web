package server

import (
    "errors"
    "fmt"
    "google.golang.org/grpc"
    "net"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "time"
    "context"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/log"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
)

var c chan os.Signal
var wg sync.WaitGroup

type IServer interface {
    init()
    Run()
}

type ServerOpt struct {
    EtcdServer           string
    Prefix               string
    ServerAddr           string
    Ctx                  context.Context
    DialTimeout          time.Duration
    DialKeepAlive        time.Duration
    NetType              string
    MaxThreadCount       string
    RegisterServerFunc   RegisterServer
    ServiceStruct        interface{}
}

type Server struct {
    opts                  ServerOpt
    listenConnector       net.Listener
    registrar             *etcdv3.Registrar
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
        EtcdServer: "127.0.0.1:2379",
        Prefix: "/services/book/",
        ServerAddr: ":0",
        Ctx: context.Background(),
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
        MaxThreadCount: "1024*1024",
        NetType: "tcp",
        RegisterServerFunc: nil,
        ServiceStruct: nil,
    }

    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func (s *Server)init(){
    options := etcdv3.ClientOptions{
        DialTimeout: s.opts.DialTimeout,
        DialKeepAlive: s.opts.DialKeepAlive,
    }
    //创建etcd连接
    client, err := etcdv3.NewClient(s.opts.Ctx, []string{s.opts.EtcdServer}, options)
    if err != nil {
        panic(err)
    }

    ls, _ := net.Listen("tcp", s.opts.ServerAddr)

    port := ls.Addr().(*net.TCPAddr).Port
    s.listenConnector = ls
    instance := ":" + strconv.Itoa(port)
    // 创建注册器
    registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
        Key:   s.opts.Prefix + instance,
        Value: instance,
    }, log.NewNopLogger())

    // 注册器启动注册
    registrar.Register()
    s.registrar = registrar

    client.SetKV(s.opts.Prefix + "thread", s.opts.MaxThreadCount)
}

func (s *Server)Run(){
    if s.opts.RegisterServerFunc == nil || s.opts.ServiceStruct == nil {
        panic(errors.New("RegisterServerFunc and ServiceStruct must set"))
    }
    c = make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    wg.Add(1)

    gs := grpc.NewServer(grpc.UnaryInterceptor(grpc_transport.Interceptor))
    defer func(){
        s.registrar.Deregister()
        s.registrar = nil
        gs.Stop()
    }()

    s.opts.RegisterServerFunc(gs, s.opts.ServiceStruct)
    go gs.Serve(s.listenConnector)
    go Producer()
    wg.Wait()
}

func Producer(){
    Loop:
    for{
        select {
        case s := <-c:
            fmt.Println()
            fmt.Println("Producer | get", s)
            break Loop
        default:
        }
        time.Sleep(500 * time.Millisecond)
    }
    wg.Done()
}

func DefaultdecodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}

func DefaultencodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}