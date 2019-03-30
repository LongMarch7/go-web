package server

import (
    "encoding/json"
    "github.com/go-kit/kit/endpoint"
    "net/http"
    "github.com/gorilla/mux"
    httptransport "github.com/go-kit/kit/transport/http"
    stdprometheus "github.com/prometheus/client_golang/prometheus/promhttp"
    "context"
)
type HealthRequest struct{

}
type errorer interface {
    error() error
}
// decode health check
func defaultDecodeRequest(_ context.Context, _ *http.Request) (interface{}, error) {
    return HealthRequest{}, nil
}

func defaultEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    if e, ok := response.(errorer); ok && e.error() != nil {
        // Not a Go kit transport error, but a business-logic error.
        // Provide those as HTTP errors.
        encodeError(ctx, e.error(), w)
        return nil
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    return json.NewEncoder(w).Encode(response)
}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
    if err == nil {
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": err.Error(),
    })
}

type Endpoints struct {
    HealthEndpoint endpoint.Endpoint
}
type HealthService interface{
    HealthCheck() bool
}

type DefaultHealthService struct {}
func (s* DefaultHealthService) HealthCheck() bool{
    return true
}
type HealthResponse struct {
    Status bool
}
func MakeHealthEndpoint(svc HealthService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        status := svc.HealthCheck()
        return HealthResponse{Status: status }, nil
    }
}
func MakeHttpHandler(_ context.Context, endpoint Endpoints) http.Handler {
    r := mux.NewRouter()
    options := []httptransport.ServerOption{
        httptransport.ServerErrorEncoder(encodeError),
    }

    //GET /health
    r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
        endpoint.HealthEndpoint,
        defaultDecodeRequest,
        defaultEncodeResponse,
        options...,
    ))

    // GET /metrics
    r.Path("/metrics").Handler(stdprometheus.Handler())
    return r
}
