package prometheus

import (
    "github.com/go-kit/kit/endpoint"
    kitprometheus "github.com/go-kit/kit/metrics/prometheus"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type PrometheusFunc interface{
    PrometheusCounterMiddleware( string, []string, float64) endpoint.Middleware
    PrometheusSummaryMiddleware( string, []string) endpoint.Middleware
    PrometheusHistogramMiddleware( string,  []string) endpoint.Middleware
    PrometheusGaugeMiddleware( string,  []string, float64, bool) endpoint.Middleware
}
type prometheusServer struct{
    Counter map[string]*kitprometheus.Counter
    Summary map[string]*kitprometheus.Summary
    Histogram map[string]*kitprometheus.Histogram
    Gauge map[string]*kitprometheus.Gauge
}

const (
    Counter_TYPE      int8 = 0
    Summary_TYPE      int8 = 1
    Gauge_TYPE        int8 = 2
    Histogram_TYPE    int8 = 3
)

type PrometheusConfig struct{
    FieldKeys []string
    Namespace string
    Subsystem string
    Name      string
    Help      string
    Type      int8
    Buckets []float64
}

func NewPrometheus(config []PrometheusConfig) PrometheusFunc{

    var server prometheusServer
    server.Counter = make(map[string]*kitprometheus.Counter)
    server.Summary = make(map[string]*kitprometheus.Summary)
    server.Histogram = make(map[string]*kitprometheus.Histogram)
    server.Gauge = make(map[string]*kitprometheus.Gauge)
    for _, value := range config{
        switch value.Type{
        case Counter_TYPE:
            requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
                Namespace: value.Namespace,
                Subsystem: value.Subsystem,
                Name:      value.Name + "_Counter",
                Help:      value.Help,
            }, value.FieldKeys)
            server.Counter[value.Name] = requestCount
        case Summary_TYPE:
            requestSummary := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
                Namespace: value.Namespace,
                Subsystem: value.Subsystem,
                Name:      value.Name + "_Summary",
                Help:      value.Help,
            }, value.FieldKeys)
            server.Summary[value.Name] = requestSummary
        case Gauge_TYPE:
            requestGauge := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
                Namespace: value.Namespace,
                Subsystem: value.Subsystem,
                Name:      value.Name + "_Gauge",
                Help:      value.Help,
            }, value.FieldKeys)
            server.Gauge[value.Name] = requestGauge
        case Histogram_TYPE:
            requestHistogram := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
                Namespace: value.Namespace,
                Subsystem: value.Subsystem,
                Name:      value.Name + "_Histogram",
                Help:      value.Help,
                Buckets: value.Buckets,
            }, value.FieldKeys)
            server.Histogram[value.Name] = requestHistogram
        }
    }
    return &server
}
