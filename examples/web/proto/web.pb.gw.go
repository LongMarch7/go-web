// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: src/github.com/LongMarch7/go-web/examples/web/proto/web.proto

/*
Package proto is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package proto

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/LongMarch7/go-web/runtime/base"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ grpclog.Logger

var (
	filter_BookService_GetBookInfo_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_BookService_GetBookInfo_0(ctx context.Context, marshaler runtime.Marshaler, client BookServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Message
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_BookService_GetBookInfo_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetBookInfo(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterBookServiceHandlerClient registers the http handlers for service BookService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "BookServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "BookServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "BookServiceClient" to call the correct interceptors.
type RequestFunc func(context.Context, runtime.Marshaler, BookServiceClient, *http.Request, map[string]string) (proto.Message, runtime.ServerMetadata, error)

func RegisterBookServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, endpoint endpoint.Endpoint, extend interface{}) error {
	//var commonFunc = func(w http.ResponseWriter, req *http.Request, pathParams map[string]string, request RequestFunc, conn *grpc.ClientConn) (runtime.Marshaler, proto.Message, error){
	//        var client = NewBookServiceClient(conn)
	//
	//		ctx, cancel := context.WithCancel(req.Context())
	//
	//		defer cancel()
	//		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	//		rctx, err := runtime.AnnotateContext(ctx, mux, req)
	//		if err != nil {
	//			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
	//			return outboundMarshaler, nil, errors.New("runtime.AnnotateContext error")
	//		}
	//		resp, md, err := request(rctx, inboundMarshaler, client, req, pathParams)
	//		ctx = runtime.NewServerMetadataContext(ctx, md)
	//		if err != nil {
	//			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
	//			return outboundMarshaler, nil, errors.New("runtime.NewServerMetadataContext error")
	//		}
	//       return outboundMarshaler, resp, nil
	//   }

	mux.Handle("GET", pattern_BookService_GetBookInfo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		manager := base.NewManager(extend)
		manager.Method = "GET"
		manager.Url = req.URL.Path
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		manager.Handler = func(rctx context.Context, conn *grpc.ClientConn) error {
			var client = NewBookServiceClient(conn)
			resp, md, err := request_BookService_GetBookInfo_0(rctx, inboundMarshaler, client, req, pathParams)
			ctx = runtime.NewServerMetadataContext(ctx, md)
			if err != nil {
				runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
				return errors.New("runtime.NewServerMetadataContext error")
			}
			if err != nil {
				return err
			}

			forward_BookService_GetBookInfo_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

			return nil
		}
		endpoint(rctx, manager)
	})

	return nil
}

var (
	pattern_BookService_GetBookInfo_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"application", "index"}, ""))
)

var (
	forward_BookService_GetBookInfo_0 = runtime.ForwardResponseMessage
)
