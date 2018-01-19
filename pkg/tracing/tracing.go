package tracing

import (
	"fmt"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
)

// https://github.com/jaegertracing/jaeger-client-go/blob/master/config/config.go

// Init returns an instance of Jaeger Tracer
func Init(serviceName string, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer) {
	// create configuration
	cfg := &config.Configuration{
		// Valid values for Param field are:
		// - for "const" sampler, 0 or 1 for always false/true respectively
		// - for "probabilistic" sampler, a probability between 0 and 1
		// - for "rateLimiting" sampler, the number of spans per second
		// - for "remote" sampler, param is the same as for "probabilistic"
		//   and indicates the initial sampling rate before the actual one
		//   is received from the mothership
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false, // log all spans to stdout for debug purposes
			BufferFlushInterval: 1 * time.Second,
		},
	}

	// instantiate tracer
	tracer, closer, err := cfg.New(
		serviceName,
		config.Logger(jaeger.StdLogger),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}

// Another feature that built into the the client libraries was the ability
// to poll the tracing backend for the sampling strategy. When a service
// receives a request that has no tracing metadata, the tracing instrumentation
// usually starts a new trace for that request by generating a new random trace
// ID. However, most production tracing systems, do not profile every single
// trace or record it in storage. Doing so would create a prohibitively large
// volume of traffic from the services to the tracing backend, possibly
// orders of magnitude larger than the actual business traffic handled by
// the services. Instead, most tracing systems sample only a small percentage
// of traces and only profile and record those sampled traces. The exact
// algorithm for making a sampling decision is what we call a sampling strategy.
