package plugin

import (
	"github.com/LongMarch7/go-web/transport/server"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/ratelimit"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"golang.org/x/time/rate"
	"time"
	"github.com/LongMarch7/go-web/plugin/logger-middleware"
	"github.com/LongMarch7/go-web/plugin/instrumenting-middleware"
	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"context"
)

type EnAndDeFun func(_ context.Context, req interface{}) (interface{}, error)

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
	logLogger                     log.Logger
	duration                      metrics.Histogram
	otTracer                      stdopentracing.Tracer
	zipkinTracer                  *stdzipkin.Tracer
	makeEndpoint               func() endpoint.Endpoint
	decodeFun                     grpc_transport.DecodeRequestFunc
	encodeFun                     grpc_transport.EncodeResponseFunc
}

type Plugin struct {
	opts                  PluginOpt
	server *grpc_transport.Server
}

type POption func(o *PluginOpt)
func NewPlugin(opts ...POption) *grpc_transport.Server {
	return newPlugin(opts...).server
}

func newPlugin(opts ...POption) *Plugin {
	options := newOptions(opts...)

	s := &Plugin{
		opts: options,
		server: nil,
	}
	var pluginEndpoint endpoint.Endpoint
	if options.makeEndpoint != nil {
		pluginEndpoint = options.makeEndpoint()
		pluginEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(options.ratelimitInterval), options.ratelimitCap))(pluginEndpoint)
		pluginEndpoint = newHystrixbyDefaultConfig(options)(pluginEndpoint)
		pluginEndpoint = kitopentracing.TraceServer(options.otTracer, options.prefix+"---"+options.methodName)(pluginEndpoint)
		pluginEndpoint = kitzipkin.TraceEndpoint(options.zipkinTracer, options.prefix+"---"+options.methodName)(pluginEndpoint)
		pluginEndpoint = loggerMid.LoggingMiddleware(log.With(options.logLogger, "method", options.methodName))(pluginEndpoint)
		pluginEndpoint = instrumentingMid.InstrumentingMiddleware(options.duration.With("method", options.methodName))(pluginEndpoint)
		s.server = grpc_transport.NewServer(
			pluginEndpoint,
			options.decodeFun,
			options.encodeFun,
		)
	}
	return s
}

func newOptions(opts ...POption) PluginOpt {
	zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	opt := PluginOpt{
		prefix: "/gateway/default",
		methodName: "default",
		ratelimitInterval: time.Millisecond * 10,
		ratelimitCap: 100,
		hystrixTimeout: 1000,
		hystrixErrorPercentThreshold: 50,
		hystrixSleepWindow: 5000,
		hystrixMaxConcurrentRequests: 100,
		hystrixRequestVolumeThreshold: 50,
		logLogger: log.NewNopLogger(),
		duration: discard.NewHistogram(),
		otTracer: opentracing.GlobalTracer(),
		zipkinTracer: zkt,
		makeEndpoint: nil,
		decodeFun: server.DefaultdecodeRequest,
		encodeFun: server.DefaultencodeResponse,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
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