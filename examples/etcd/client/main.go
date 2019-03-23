package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/sd"
    "github.com/LongMarch7/go-web/util/sd/etcdv3"
    "github.com/go-kit/kit/sd/lb"
    "google.golang.org/grpc"
    "github.com/LongMarch7/go-web/examples/etcd/book"
    "io"
    "time"
    "github.com/LongMarch7/go-web/transport/pool"
)

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
    instancer, err := etcdv3.NewInstancer(client, *prefix, logger, pool.Update)
    if err != nil {
        panic(err)
    }
    //创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
    endpointer := sd.NewEndpointer(instancer, reqFactory, logger)
    //创建负载均衡器
    balancer := lb.NewRoundRobin(endpointer)
    fmt.Printf("book = %s", client.GetKV("book"))
    /**
    我们可以通过负载均衡器直接获取请求的endPoint，发起请求
    reqEndPoint,_ := balancer.Endpoint()
    */

    /**
    也可以通过retry定义尝试次数进行请求
    */
    reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

    //现在我们可以通过 endPoint 发起请求了
    req := struct{}{}
    //if _, err = reqEndPoint(ctx, req); err != nil {
    //    panic(err)
    //}
    for i := 1; i <= 10; i++ {
        fmt.Println("请求count: ",i)
        if _, err = reqEndPoint(ctx, req); err != nil {
            panic(err)
        }
    }
}


//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        fmt.Println("请求服务: ", instanceAddr)
        conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
        if err != nil {
            fmt.Println(err)
            panic("connect error")
        }
        defer func() {
            conn.Close()
        }()
        fmt.Println("stat=",conn.GetState().String())
        bookClient := book.NewBookServiceClient(conn)
        bi,_:=bookClient.GetBookInfo(context.Background(),&book.BookInfoParams{BookId:1})
        if bi == nil{
            return nil,errors.New("GetBookInfo error")
        }
        fmt.Println("获取书籍详情")
        fmt.Println("bookId: 1", " => ", "bookName:", bi.BookName)

        bl,_ := bookClient.GetBookList(context.Background(), &book.BookListParams{Page:1, Limit:10})
        if bl == nil{
            return nil,errors.New("GetBookList error")
        }
        fmt.Println("获取书籍列表")
        for _,b := range bl.BookList {
            fmt.Println("bookId:", b.BookId, " => ", "bookName:", b.BookName)
        }
        return nil,nil
    },nil,nil
}