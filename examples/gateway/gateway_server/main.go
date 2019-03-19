package main

import (
    "context"
    "fmt"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
    "github.com/go-kit/kit/endpoint"
    "github.com/LongMarch7/go-web/examples/gateway/book"
    "google.golang.org/grpc"
    "net"
    "github.com/go-kit/kit/sd/etcdv3"
    "github.com/go-kit/kit/log"
    "strconv"
    "time"
)

type BookServer struct {
    bookListHandler  grpc_transport.Handler
    bookInfoHandler  grpc_transport.Handler
}

//通过grpc调用GetBookInfo时,GetBookInfo只做数据透传, 调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (s *BookServer) GetBookInfo(ctx context.Context, in *book.BookInfoParams) (*book.BookInfo, error) {
    _, rsp, err := s.bookInfoHandler.ServeGRPC(ctx, in)
    if err != nil {
        return nil, err
    }
    return rsp.(*book.BookInfo),err
}

//通过grpc调用GetBookList时,GetBookList只做数据透传, 调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (s *BookServer) GetBookList(ctx context.Context, in *book.BookListParams) (*book.BookList, error) {
    _, rsp, err := s.bookListHandler.ServeGRPC(ctx, in)
    if err != nil {
        return nil, err
    }
    return rsp.(*book.BookList),err
}

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

func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}

func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}

func main() {
    var (
        //etcd服务地址
        etcdServer = "127.0.0.1:2379"
        //服务的信息目录
        prefix     = "/services/book/"
        //服务实例注册的路径
        key        = prefix
        ctx        = context.Background()
        //服务监听地址
        serviceAddress = ":12306"
    )

    //etcd的连接参数
    options := etcdv3.ClientOptions{
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
    }
    //创建etcd连接
    client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
    if err != nil {
        panic(err)
    }

    ls, _ := net.Listen("tcp", serviceAddress)

    port := ls.Addr().(*net.TCPAddr).Port
    instance := ":" + strconv.Itoa(port)
    // 创建注册器
    registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
        Key:   key+instance,
        Value: instance,
    }, log.NewNopLogger())

    // 注册器启动注册
    registrar.Register()


    bookServer := new(BookServer)
    bookListHandler := grpc_transport.NewServer(
        makeGetBookListEndpoint(),
        decodeRequest,
        encodeResponse,
    )
    bookServer.bookListHandler = bookListHandler


    bookInfoHandler := grpc_transport.NewServer(
        makeGetBookInfoEndpoint(),
        decodeRequest,
        encodeResponse,
    )
    bookServer.bookInfoHandler = bookInfoHandler

    gs := grpc.NewServer(grpc.UnaryInterceptor(grpc_transport.Interceptor))
    book.RegisterBookServiceServer(gs, bookServer)
    gs.Serve(ls)
}