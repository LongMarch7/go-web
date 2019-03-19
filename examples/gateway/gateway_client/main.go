package main

import (
    "flag"
    "net/http"

    "context"
    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "google.golang.org/grpc"

    "github.com/LongMarch7/go-web/examples/gateway/book"
)

var (
    // the go.micro.srv.greeter address
    endpoint = flag.String("endpoint", "localhost:12306", "go.micro.srv.greeter address")
)

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}

    err := book.RegisterBookServiceHandlerFromEndpoint(ctx, mux, *endpoint, opts)
    if err != nil {
        return err
    }

    return http.ListenAndServe(":8080", mux)
}

func main() {
    flag.Parse()

    defer glog.Flush()

    if err := run(); err != nil {
        glog.Fatal(err)
    }
}
