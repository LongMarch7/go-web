package main

import (
    "context"
    "flag"
    "github.com/go-kit/kit/endpoint"
    "github.com/LongMarch7/go-web/examples/web/proto"
    "github.com/LongMarch7/go-web/runtime/consul-server"
    "github.com/LongMarch7/go-web/plugin"
    "google.golang.org/grpc/grpclog"
    zapLog "github.com/LongMarch7/go-web/plugin/zap-log"
    "github.com/LongMarch7/go-web/app"
    "github.com/LongMarch7/go-web/view"
)

func makeIndexEndpoint(tpl view.View) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        //请求列表时返回 书籍列表
        //req := request.(*proto.BookListParams)
        //grpclog.Info("registerGrpc faild limit:",req.Limit)
        //grpclog.Info("page:",req.Page)
        //bl := new(proto.BookList)
        //bl.BookList = append(bl.BookList, &proto.BookInfo{BookId: 1,BookName:"21天精通php"})
        //bl.BookList = append(bl.BookList, &proto.BookInfo{BookId: 2,BookName:"21天精通java"})
        test := new(proto.Message)
        test.Text = "hello"
        grpclog.Info("test")
        return test,nil
    }
}


func main() {
    etcdServer := flag.String("e","http://localhost:8500","etcd service addr")
    prefix := flag.String("n","bookServer","prefix value")
    serviceAddress := flag.String("s","127.0.0.1","server addr")
    servicePort := flag.Int("p",0,"server port")
    threadMax := flag.String("c","1024","server thread pool max thread count")
    flag.Parse()
    ctx := context.Background()

    grpclog.SetLoggerV2(zapLog.NewDefaultLoggerConfig().NewLogger())

    Init := func(tpl view.View) interface {}{
        bookServer := new(proto.DefaultBookServiceServer)
        bookServer.GetBookInfoHandler = plugin.NewPlugin(
            plugin.Prefix(*prefix),
            plugin.MethodName("GetBookInfo"),
            plugin.MakeEndpoint(makeIndexEndpoint(tpl)))
        return bookServer
    }

    microServer :=app.NewMicroServer(
        server.ConsulAddr(*etcdServer),
        server.Prefix(*prefix),
        server.ServerAddr(*serviceAddress),
        server.ServerPort(*servicePort),
        server.Ctx(ctx),
        server.MaxThreadCount(*threadMax),                     //must set
        server.RegisterServiceFunc(proto.RegisterBookServiceServer), //must set)
    )
    microServer.AddPlugin(Init(microServer.Template))
    microServer.Run()
}