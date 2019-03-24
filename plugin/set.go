package plugin

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"golang.org/x/time/rate"
	"time"
	"github.com/LongMarch7/go-web/plugin/logger"
	"github.com/LongMarch7/go-web/plugin/instrumenting"
	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)


type PluginOpt struct {
	prefix                        string
	methodName                    string
	ratelimitInterval             time.Duration
	ratelimitCap                  int
	hystrixName                   string
	hystrixTimeout                int
	hystrixMaxConcurrentRequests  int
	hystrixRequestVolumeThreshold int
	hystrixSleepWindow            int
	hystrixErrorPercentThreshold  int
	logger                        log.Logger
	duration                      metrics.Histogram
	otTracer                      stdopentracing.Tracer
	zipkinTracer                  *stdzipkin.Tracer
	makeSumEndpoint               func() endpoint.Endpoint
}

type Plugin struct {
	opts                  PluginOpt
	pluginEndpoint endpoint.Endpoint
}

type POption func(o *PluginOpt)
func NewPlugin(opts ...POption) endpoint.Endpoint {
	return newPlugin(opts...).pluginEndpoint
}

func newPlugin(opts ...POption) *Plugin {
	options := newOptions(opts...)

	var pluginEndpoint endpoint.Endpoint
	{
		pluginEndpoint = options.makeSumEndpoint()
		pluginEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(options.ratelimitInterval), options.ratelimitCap))(pluginEndpoint)
		pluginEndpoint = newHystrixbyDefaultConfig(options)(pluginEndpoint)
		pluginEndpoint = opentracing.TraceServer(options.otTracer, options.prefix + "---" + options.methodName)(pluginEndpoint)
		pluginEndpoint = zipkin.TraceEndpoint(options.zipkinTracer, options.prefix + "---" + options.methodName)(pluginEndpoint)
		pluginEndpoint = logger.LoggingMiddleware(log.With(options.logger, "method", options.methodName))(pluginEndpoint)
		pluginEndpoint = instrumenting.InstrumentingMiddleware(options.duration.With("method", options.methodName))(pluginEndpoint)
	}
	s := &Plugin{
		opts: options,
		pluginEndpoint: pluginEndpoint,
	}
	return s
}
func newOptions(opts ...POption) PluginOpt {
	opt := PluginOpt{

	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func (p *Plugin) Run(){

}

func newHystrixbyDefaultConfig(opt PluginOpt) endpoint.Middleware{
	commandName := opt.hystrixName
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout: opt.hystrixTimeout,
		ErrorPercentThreshold: opt.hystrixErrorPercentThreshold,
		SleepWindow: opt.hystrixSleepWindow,
		MaxConcurrentRequests: opt.hystrixMaxConcurrentRequests,
		RequestVolumeThreshold: opt.hystrixRequestVolumeThreshold,
	})

	return circuitbreaker.Hystrix(commandName)
}
