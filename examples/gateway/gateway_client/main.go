package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/go-kit/kit/sd"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/sd/lb"
    "google.golang.org/grpc"
    "io"
    "net/http"
    "time"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/examples/gateway/book"
    "github.com/LongMarch7/go-web/transport/client"
    //"github.com/daviddengcn/go-villa"
)

//
//func run() error {
//    ctx := context.Background()
//    ctx, cancel := context.WithCancel(ctx)
//    defer cancel()
//
//    mux := runtime.NewServeMux()
//    opts := []grpc.DialOption{grpc.WithInsecure()}
//
//    err := book.RegisterBookServiceHandlerFromEndpoint(ctx, mux, *endpoint, opts)
//    if err != nil {
//        return err
//    }
//
//    return http.ListenAndServe(":8080", mux)
//}
//
//func main() {
//    flag.Parse()
//
//    defer glog.Flush()
//
//    if err := run(); err != nil {
//        glog.Fatal(err)
//    }
//}
func main() {
    etcdServer := flag.String("e","127.0.0.1:2379","etcd service addr")
    prefix := flag.String("p","/services/book/","prefix value")
    flag.Parse()
    ctx := context.Background()

    options := etcdv3.ClientOptions{
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
    }
    //连接注册中心
    client, err := etcdv3.NewClient(ctx, []string{*etcdServer}, options)
    if err != nil {
        panic(err)
    }
    logger := log.NewNopLogger()
    //创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
    instancer, err := etcdv3.NewInstancer(client, *prefix, logger)
    if err != nil {
        panic(err)
    }
    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, reqFactory, logger)
    //创建负载均衡器
    balancer := lb.NewRoundRobin(endpointer)

    /**
    我们可以通过负载均衡器直接获取请求的endPoint，发起请求
    reqEndPoint,_ := balancer.Endpoint()
    */

    /**
    也可以通过retry定义尝试次数进行请求
    */
    reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()

    err = book.RegisterBookServiceHandlerClient(ctx, mux, reqEndPoint, nil)
    if err != nil {
        panic(err)
    }

    http.ListenAndServe(":8080", mux)
}

//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        fmt.Println("请求服务: ", instanceAddr)
        manage :=request.(*client.BaseGatewayManager)
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