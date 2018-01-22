package tracing

import (
	"testing"
	"time"

	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-lib/metrics/go-kit"
	"github.com/uber/jaeger-lib/metrics/go-kit/expvar"
)

func TestInit(t *testing.T) {

	metricsFactory := xkit.Wrap("test", expvar.NewFactory(10))
	tracer, closer, err := Init(
		"test",
		metricsFactory.Namespace(info.Report.Program, nil),
	)
	if err != nil {
		t.Error("Could not initialize tracer")
	}

	// start a new span
	span := tracer.StartSpan("test")

	// set a new tag in the span
	span.SetTag("hello", "test")

	// Set baggage in the span
	span.SetBaggageItem("baggage", "test")

	// wait
	time.Sleep(100 * time.Millisecond)

	// Log to the span
	span.LogFields(
		log.String("event", "test complete"),
	)

	// close things down
	span.Finish()
	closer.Close()
}
