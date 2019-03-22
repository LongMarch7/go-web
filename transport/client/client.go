package client

import (
    "context"
    "fmt"
    "github.com/go-kit/kit/sd"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/sd/lb"
    "google.golang.org/grpc"
    "io"
    "strconv"
    "time"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type GatewayHandler func(conn *grpc.ClientConn)
type BaseGatewayManager struct {
    Handler GatewayHandler
    Extend  interface{}
}

func NewManager(extend interface{}) *BaseGatewayManager {
    return &BaseGatewayManager{
        Extend: extend,
    }
}

type RegisterHandlerClient func( context.Context, *runtime.ServeMux, endpoint.Endpoint, interface{}) error
type Client interface {
    Register(opt ClientOpt, register RegisterHandlerClient)
    DeRegister()
}
type ClientOpt struct {
    EtcdServer      string
    Prefix          string
    Mux             *runtime.ServeMux
    Ctx             context.Context
    DialTimeout     time.Duration
    DialKeepAlive   time.Duration
    Cancel          context.CancelFunc
    RegisterGrpc    RegisterHandlerClient
    Factory         sd.Factory
    RetryTime       time.Duration
    RetryCount      int
    MaxThreadCount  int
}

func NewClientOpt(etcd string, prefix string , mux *runtime.ServeMux, ctx context.Context, handler RegisterHandlerClient) *ClientOpt{
    return &ClientOpt{
        EtcdServer: etcd,
        Prefix: prefix,
        Mux: mux,
        Ctx: ctx,
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
        Factory: defaultReqFactory,
        RetryTime: time.Second * 3,
        RetryCount: 3,
        RegisterGrpc: handler,
    }
}

func (c *ClientOpt)Register(){
    options := etcdv3.ClientOptions{
        DialTimeout: c.DialTimeout,
        DialKeepAlive: c.DialKeepAlive,
    }
    //连接注册中心
    client, err := etcdv3.NewClient(c.Ctx, []string{c.EtcdServer}, options)
    if err != nil {
        panic(err)
    }
    logger := log.NewNopLogger()
    //创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
    instancer, err := etcdv3.NewInstancer(client, c.Prefix, logger)
    if err != nil {
        panic(err)
    }
    c.MaxThreadCount, _ =  strconv.Atoi(string(client.GetKV(c.Prefix + "thread")))
    if c.MaxThreadCount <= 0 {
        c.MaxThreadCount = 50
    }
    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, c.Factory, logger)
    //创建负载均衡器
    balancer := lb.NewRoundRobin(endpointer)

    /**
    我们可以通过负载均衡器直接获取请求的endPoint，发起请求
    reqEndPoint,_ := balancer.Endpoint()
    */

    /**
    也可以通过retry定义尝试次数进行请求
    */
    reqEndPoint := lb.Retry(c.RetryCount, c.RetryTime, balancer)

    ctx, cancel := context.WithCancel(c.Ctx)
    c.Cancel = cancel

    err = c.RegisterGrpc(ctx, c.Mux, reqEndPoint, nil)
    if err != nil {
        panic(err)
    }
}

func (c *ClientOpt)DeRegister(){
    c.Cancel()
}

func defaultReqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        fmt.Println("请求服务: ", instanceAddr)
        manage :=request.(*BaseGatewayManager)
        conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
        if err != nil {
            fmt.Println(err)
            panic("connect error")
        }
        defer func() {
            conn.Close()
        }()
        manage.Handler(conn)
        return nil,nil
    },nil,nil
}
