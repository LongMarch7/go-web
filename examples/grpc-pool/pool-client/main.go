package main

import (
    "context"
    "flag"
    "net/http"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/examples/grpc-pool/book"
    "github.com/LongMarch7/go-web/transport/client"
    "time"
)

func main() {
    etcdServer := flag.String("e","127.0.0.1:2379","etcd service addr")
    prefix := flag.String("p","/services/book/","prefix value")
    flag.Parse()
    ctx := context.Background()
    mux := runtime.NewServeMux()

    client1 := client.NewClient(
        client.EtcdServer(*etcdServer),
        client.Prefix(*prefix),
        client.Mux(mux),                                               //must set
        client.Ctx(ctx),
        client.RegisterGrpc(book.RegisterBookServiceHandlerClient),   //must se
        client.RetryCount(3),
        client.RetryTime(time.Second * 3),
    )
    defer client1.DeRegister()

    http.ListenAndServe(":8080", mux)
}