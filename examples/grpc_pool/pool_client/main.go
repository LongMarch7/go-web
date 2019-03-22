package main

import (
    "context"
    "flag"
    "net/http"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/examples/grpc_pool/book"
    "github.com/LongMarch7/go-web/transport/client"
)

func main() {
    etcdServer := flag.String("e","127.0.0.1:2379","etcd service addr")
    prefix := flag.String("p","/services/book/","prefix value")
    flag.Parse()
    ctx := context.Background()
    mux := runtime.NewServeMux()

    client1 := client.NewClientOpt(*etcdServer, *prefix, mux, ctx, book.RegisterBookServiceHandlerClient)
    client1.Register()
    defer client1.DeRegister()

    http.ListenAndServe(":8080", mux)
}