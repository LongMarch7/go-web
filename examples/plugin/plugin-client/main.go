package main

import (
    "context"
    "flag"
    "fmt"
    "net/http"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/LongMarch7/go-web/examples/plugin/book"
    "github.com/LongMarch7/go-web/transport/client"
    "os"
    "os/signal"
    "sync"
    "time"
)

var c chan os.Signal
var wg sync.WaitGroup

func Producer(){
Loop:
    for{
        select {
        case s := <-c:
            fmt.Println()
            fmt.Println("Producer | get", s)
            break Loop
        default:
        }
        time.Sleep(500 * time.Millisecond)
    }
    wg.Done()
}

func main() {
    etcdServer := flag.String("e","127.0.0.1:2379","etcd service addr")
    prefix := flag.String("p","/services/book/","prefix value")
    flag.Parse()

    ctx := context.Background()
    mux := runtime.NewServeMux()
    defer func(){
        mux = nil
    }()

    c = make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    wg.Add(1)

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

    go http.ListenAndServe(":8080", mux)
    go Producer()
    wg.Wait()
}

