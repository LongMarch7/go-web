package server

import (
    "net"
    "strconv"
    "time"
    "context"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/log"
)


type IServer interface {
    init()
    GetListener() net.Listener
}

type ServerOpt struct {
    EtcdServer      string
    Prefix          string
    ServerAddr        string
    Ctx             context.Context
    DialTimeout     time.Duration
    DialKeepAlive   time.Duration
    NetType         string
    MaxThreadCount  string
}

type Server struct {
    opts ServerOpt
    listenConnector net.Listener
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

    client.SetKV(s.opts.Prefix + "thread", s.opts.MaxThreadCount)
}

func (s *Server)GetListener() net.Listener{
    return s.listenConnector
}

func DefaultdecodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}

func DefaultencodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}