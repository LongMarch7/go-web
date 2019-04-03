package prometheus

import (
    "github.com/go-kit/kit/endpoint"
    "context"
    "time"
)

func (p *prometheusServer)PrometheusCounterMiddleware(name string, lvs []string, count float64) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        action, ok := p.Counter[name + "_Counter"]
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
        action, ok := p.Summary[name + "_Summary"]
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
        action, ok := p.Histogram[name + "_Histogram"]
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
        action, ok := p.Gauge[name + "_Gauge"]
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
