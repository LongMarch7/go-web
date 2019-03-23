package client

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-kit/kit/sd"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/sd/lb"
    "google.golang.org/grpc"
    "io"
    "time"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/transport/pool"
)

type GatewayHandler func(conn *grpc.ClientConn)
type BaseGatewayManager struct {
    Handler GatewayHandler
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
    EtcdServer      string
    Prefix          string
    Mux             *runtime.ServeMux
    Ctx             context.Context
    DialTimeout     time.Duration
    DialKeepAlive   time.Duration
    RegisterGrpc    RegisterHandlerClient
    Factory         sd.Factory
    RetryTime       time.Duration
    RetryCount      int
    Extend          interface{}
    Manager         interface{}
}

type Client struct {
    opts      ClientOpt
    cancel    context.CancelFunc
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
        EtcdServer: "127.0.0.1:2379",
        Prefix: "/services/book/",
        Ctx: context.Background(),
        Mux: nil,
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
        Factory: defaultReqFactory,
        RetryTime: time.Second * 3,
        RetryCount: 3,
        RegisterGrpc: nil,
        Extend: nil,
        Manager: nil,
    }

    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func (c *Client)init(){
    if c.opts.Mux == nil || c.opts.RegisterGrpc == nil{
        fmt.Println("mux and grpc need set")
        return
    }
    options := etcdv3.ClientOptions{
        DialTimeout: c.opts.DialTimeout,
        DialKeepAlive: c.opts.DialKeepAlive,
    }
    //连接注册中心
    client, err := etcdv3.NewClient(c.opts.Ctx, []string{c.opts.EtcdServer}, options)
    if err != nil {
        panic(err)
    }
    logger := log.NewNopLogger()
    //创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
    instancer, err := etcdv3.NewInstancer(client, c.opts.Prefix, logger, pool.Update)
    if err != nil {
        panic(err)
    }

    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, c.opts.Factory, logger)
    //创建负载均衡器
    balancer := lb.NewRoundRobin(endpointer)

    /**
    我们可以通过负载均衡器直接获取请求的endPoint，发起请求
    reqEndPoint,_ := balancer.Endpoint()
    */

    /**
    也可以通过retry定义尝试次数进行请求
    */
    reqEndPoint := lb.Retry(c.opts.RetryCount, c.opts.RetryTime, balancer)

    ctx, cancel := context.WithCancel(c.opts.Ctx)
    c.cancel = cancel

    if c.opts.Manager == nil {
        c.opts.Manager = &GrpcPoolManager{
            Opt: grpc.WithInsecure(),
            Extend: c.opts.Extend,
        }
    }
    err = c.opts.RegisterGrpc(ctx, c.opts.Mux, reqEndPoint, c.opts.Manager)
    if err != nil {
        panic(err)
    }
}

func (c *Client)DeRegister(){
    c.cancel()
    pool.Destroy()
}

func getConnectFromPool(addr string, p pool.Pool, opt grpc.DialOption) (*grpc.ClientConn, error){
    var err error = nil
    conn, ok, _ := p.Queue.Get()
    if !ok {
        time.Sleep(time.Microsecond * 100)
        conn, ok, _ =  p.Queue.Get()
    }
    if !ok{
        conn, err = grpc.Dial(addr, opt)
    }
    return conn.(*grpc.ClientConn), err
}

func putConnectToPool(conn *grpc.ClientConn, p pool.Pool) {
    ok, _ := p.Queue.Put(conn)
    if !ok {
        time.Sleep(time.Microsecond)
        ok, _ =  p.Queue.Put(&conn)
    }
    if !ok {
        conn.Close()
        conn = nil
    }
}

func defaultReqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        base :=request.(*BaseGatewayManager)
        poolManage,ok := pool.GetConnect(instanceAddr)
        if ! ok {
            return nil,errors.New("poolManage not found")
        }

        conn, err := getConnectFromPool(instanceAddr, poolManage, base.manager.(*GrpcPoolManager).Opt)
        if err != nil {
            return nil,err
        }
        defer func() {
            putConnectToPool(conn, poolManage)
        }()
        base.Handler(conn)
        return nil,nil
    },nil,nil
}