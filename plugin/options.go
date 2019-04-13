package plugin

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"time"
	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
)

func Prefix(prefix string) POption {
	return func(o *PluginOpt) {
		o.prefix = prefix
	}
}

func MethodName(methodName string) POption {
	return func(o *PluginOpt) {
		o.methodName = methodName
	}
}

func RatelimitInterval(ratelimitInterval time.Duration) POption {
	return func(o *PluginOpt) {
		o.ratelimitInterval = ratelimitInterval
	}
}

func RatelimitCap(ratelimitCap int) POption {
	return func(o *PluginOpt) {
		o.ratelimitCap = ratelimitCap
	}
}

func HystrixName(hystrixName string) POption {
	return func(o *PluginOpt) {
		o.hystrixName = hystrixName
	}
}

func HystrixTimeout(hystrixTimeout  int) POption {
	return func(o *PluginOpt) {
		o.hystrixTimeout = hystrixTimeout
	}
}

func HystrixErrorPercentThreshold(hystrixErrorPercentThreshold  int) POption {
	return func(o *PluginOpt) {
		o.hystrixErrorPercentThreshold = hystrixErrorPercentThreshold
	}
}

func HystrixSleepWindow(hystrixSleepWindow  int) POption {
	return func(o *PluginOpt) {
		o.hystrixSleepWindow = hystrixSleepWindow
	}
}

func HystrixMaxConcurrentRequests(hystrixMaxConcurrentRequests  int) POption {
	return func(o *PluginOpt) {
		o.hystrixMaxConcurrentRequests = hystrixMaxConcurrentRequests
	}
}

func HystrixRequestVolumeThreshold(hystrixRequestVolumeThreshold  int) POption {
	return func(o *PluginOpt) {
		o.hystrixRequestVolumeThreshold = hystrixRequestVolumeThreshold
	}
}

func Logger(logger  log.Logger) POption {
	return func(o *PluginOpt) {
		o.logLogger = logger
	}
}

func Duration(duration  metrics.Histogram) POption {
	return func(o *PluginOpt) {
		o.duration = duration
	}
}

func OtTracer(otTracer   stdopentracing.Tracer) POption {
	return func(o *PluginOpt) {
		o.otTracer = otTracer
	}
}

func ZipkinTracer(zipkinTracer  *stdzipkin.Tracer) POption {
	return func(o *PluginOpt) {
		o.zipkinTracer = zipkinTracer
	}
}

func MakeEndpoint(endpoint  endpoint.Endpoint) POption {
	return func(o *PluginOpt) {
		o.Endpoint = endpoint
	}
}

func DecodeFun(decodeFun  grpc_transport.DecodeRequestFunc) POption {
	return func(o *PluginOpt) {
		o.decodeFun = decodeFun
	}
}

func EncodeFun(encodeFun   grpc_transport.EncodeResponseFunc) POption {
	return func(o *PluginOpt) {
		o.encodeFun = encodeFun
	}
}