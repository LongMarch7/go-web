package zipkin

import (
    "google.golang.org/grpc/grpclog"
    zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
    stdopentracing "github.com/opentracing/opentracing-go"
)


var collector zipkinot.Collector = nil
var zip stdopentracing.Tracer = nil
func NewZipkinTracer(name string) (stdopentracing.Tracer, zipkinot.Collector){
    local_collector, err := zipkinot.NewHTTPCollector("http://127.0.0.1:9411/api/v1/spans")
    if err != nil {
        grpclog.Error("zipkinot.NewHTTPCollector faild")
        return nil,nil
    }
    var (
        debug       = false
        hostPort    = "localhost:0"
        serviceName = name
    )
    recorder := zipkinot.NewRecorder(local_collector, debug, hostPort, serviceName)
    newTracer, err := zipkinot.NewTracer(recorder)
    if err != nil {
        grpclog.Error("zipkinot.NewTracer faild")
        return nil, local_collector
    }
    return newTracer,local_collector
}

func GetZipkinTracer(name string) (stdopentracing.Tracer, zipkinot.Collector){
    if zip == nil || collector == nil{
        zip, collector = NewZipkinTracer(name)
    }
    return zip, collector
}

func CloseZipkinCollector(){
    if collector != nil {
        collector.Close()
        collector = nil
    }
}
