package client

import (
	"errors"
	"github.com/LongMarch7/go-web/runtime/pool"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	"io"
	"time"
	"context"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	pluginzipkin "github.com/LongMarch7/go-web/plugin/zipkin"
)

var InvalidateTimeout = time.Minute * 3

func getConnectFromPool(addr string, p *pool.Pool, opt ...grpc.DialOption) (*pool.ConnectManager, error){
	cManager, ok, _ := p.Queue.Get()
	if !ok {
		time.Sleep(time.Microsecond * 100)
		cManager, ok, _ =  p.Queue.Get()
	}
	if !ok{
		conn, err := grpc.Dial(addr, opt...)
		if err == nil {
			cManager = new(pool.ConnectManager)
			cManager.(*pool.ConnectManager).Conn = conn
			cManager.(*pool.ConnectManager).InvalidateDeadline = time.Now().Add(InvalidateTimeout)
		}
		return cManager.(*pool.ConnectManager), err
	}
	return cManager.(*pool.ConnectManager), nil
}

func putConnectToPool(manager *pool.ConnectManager, p *pool.Pool) {
	var ok = true

	defer func(){
		if !ok {
			manager.Conn.Close()
			manager.Conn = nil
			manager = nil
		}
	}()
	if time.Now().After(manager.InvalidateDeadline) {
		ok = false
		return
	}
	manager.InvalidateDeadline = time.Now().Add(InvalidateTimeout)
	ok, _ = p.Queue.Put(manager)
	if !ok {
		time.Sleep(time.Microsecond)
		ok, _ =  p.Queue.Put(manager)
	}

}

func defaultReqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	endpointFunc := func(ctx context.Context, request interface{}) (interface{}, error) {
		base :=request.(*BaseGatewayManager)
		poolManage,ok := pool.GetConnect(instanceAddr)
		if ! ok {
		return nil,errors.New("poolManage not found")
		}

		cManager, err := getConnectFromPool(instanceAddr, poolManage, base.manager.(*GrpcPoolManager).Opt...)
		if err != nil {
		return nil,err
		}
		defer func() {
		putConnectToPool(cManager, poolManage)
		}()
		err = base.Handler(ctx, cManager.Conn)
		return nil,err
	}
	endpointFunc = breakerMw(endpointFunc)
	zip , _ := pluginzipkin.GetZipkinTracer("gateway")
	if zip != nil {
		endpointFunc = kitopentracing.TraceClient(zip, "httpRequest")(endpointFunc)
	}
	return endpointFunc,nil,nil
}
