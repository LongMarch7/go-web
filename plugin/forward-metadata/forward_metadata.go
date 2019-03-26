package forward_metadata
import (
    "net/http"
    "golang.org/x/net/context"
    "google.golang.org/grpc/metadata"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
    prefixTracerState  = "x-b3-"
    zipkinTraceID      = prefixTracerState + "traceid"
    zipkinSpanID       = prefixTracerState + "spanid"
    zipkinParentSpanID = prefixTracerState + "parentspanid"
    zipkinSampled      = prefixTracerState + "sampled"
    zipkinFlags        = prefixTracerState + "flags"
)

var otHeaders = []string{
    zipkinTraceID,
    zipkinSpanID,
    zipkinParentSpanID,
    zipkinSampled,
    zipkinFlags}

func injectHeadersIntoMetadata(ctx context.Context, req *http.Request) metadata.MD {
    pairs := []string{}
    for _, h := range otHeaders {
        if v := req.Header.Get(h); len(v) > 0 {
            pairs = append(pairs, h, v)
        }
    }
    return metadata.Pairs(pairs...)
}

type annotator func(context.Context, *http.Request) metadata.MD

func chainGrpcAnnotators(annotators ...annotator) annotator {
    return func(c context.Context, r *http.Request) metadata.MD {
        mds := []metadata.MD{}
        for _, a := range annotators {
            mds = append(mds, a(c, r))
        }
        return metadata.Join(mds...)
    }
}

func NewServeMux() *runtime.ServeMux{
    annotators := []annotator{injectHeadersIntoMetadata}

    gwmux := runtime.NewServeMux(
        runtime.WithMetadata(chainGrpcAnnotators(annotators...)),
    )
    return gwmux
}