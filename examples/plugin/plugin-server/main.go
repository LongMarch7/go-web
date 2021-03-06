package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/go-kit/kit/endpoint"
    "github.com/LongMarch7/go-web/examples/plugin/book"
    "github.com/LongMarch7/go-web/runtime/etcd-server"
    "github.com/LongMarch7/go-web/plugin"
    "github.com/go-kit/kit/log"
    "os"
)

//创建bookList的EndPoint
func makeGetBookListEndpoint() endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        //请求列表时返回 书籍列表
        req := request.(*book.BookListParams)
        fmt.Println("limit:",req.Limit)
        fmt.Println("page:",req.Page)
        bl := new(book.BookList)
        bl.BookList = append(bl.BookList, &book.BookInfo{BookId:1,BookName:"21天精通php"})
        bl.BookList = append(bl.BookList, &book.BookInfo{BookId:2,BookName:"21天精通java"})
        return bl,nil
    }
}

//创建bookInfo的EndPoint
func makeGetBookInfoEndpoint() endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        //请求详情时返回 书籍信息
        req := request.(*book.BookInfoParams)
        fmt.Println("bookID:",req.BookId)
        b := new(book.BookInfo)
        b.BookId = req.BookId
        b.BookName = "21天精通php"
        return b,nil
    }
}


func main() {
    etcdServer := flag.String("e","127.0.0.1:2379","etcd service addr")
    prefix := flag.String("p","/services/book/","prefix value")
    serviceAddress := flag.String("s",":0","server addr")
    threadMax := flag.String("c","1024","server thread pool max thread count")
    flag.Parse()
    ctx := context.Background()

    Init := func() interface {}{
        bookServer := new(book.DefaultBookServiceServer)
        bookServer.GetBookListHandler = plugin.NewPlugin(
            plugin.Prefix(*prefix),
            plugin.Logger(log.NewLogfmtLogger(os.Stderr)),
            plugin.MethodName("GetBookListEndpoint"),
            plugin.MakeEndpoint(makeGetBookListEndpoint()))

        bookServer.GetBookInfoHandler = plugin.NewPlugin(
            plugin.Prefix(*prefix),
            plugin.Logger(log.NewLogfmtLogger(os.Stderr)),
            plugin.MethodName("GetBookInfoEndpoint"),
            plugin.MakeEndpoint(makeGetBookInfoEndpoint()))
        return bookServer
    }

    server1 := server.NewServer(
            server.EtcdServer(*etcdServer),
            server.Prefix(*prefix),
            server.ServerAddr(*serviceAddress),
            server.Ctx(ctx),
            server.MaxThreadCount(*threadMax),
            server.ServiceInit(Init),                            //must set
            server.RegisterServiceFunc(book.RegisterBookServiceServer),  //must set
        )
    server1.Run()
}