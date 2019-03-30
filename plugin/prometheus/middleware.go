package prometheus

import (
    "github.com/go-kit/kit/endpoint"
    "context"
    "time"
)

func (p *prometheusServer)PrometheusCounterMiddleware(name string, lvs []string, count float64) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        action, ok := p.Counter[name]
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if ok {
                action.With(lvs...).Add(count)
            }
            return next(ctx, request)

        }
    }
}

func (p *prometheusServer)PrometheusSummaryMiddleware(name string, lvs []string) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        action, ok := p.Summary[name]
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if ok {
                defer func(begin time.Time) {
                    action.With(lvs...).Observe(time.Since(begin).Seconds())
                }(time.Now())
            }
            return next(ctx, request)

        }
    }
}

func (p *prometheusServer)PrometheusHistogramMiddleware(name string, lvs []string) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        action, ok := p.Histogram[name]
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if ok {
                defer func(begin time.Time) {
                    action.With(lvs...).Observe(time.Since(begin).Seconds())
                }(time.Now())
            }
            return next(ctx, request)

        }
    }
}

func (p *prometheusServer)PrometheusGaugeMiddleware(name string, lvs []string, count float64, isSet bool) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        action, ok := p.Gauge[name]
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if ok {
                if isSet {
                    action.With(lvs...).Set(count)
                }else{
                    action.With(lvs...).Add(count)
                }
            }
            return next(ctx, request)
        }
    }
}
