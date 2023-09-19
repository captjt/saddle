package telemetry

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type exporter int

const (
	Jaeger exporter = iota
	CloudTrace
	StdOut
)

func NewExporter(
	exporter exporter, serviceName, environment string, sampleRate float64, endpoint ...string,
) (*trace.TracerProvider, error) {
	var (
		err error
		exp trace.SpanExporter
	)

	hn, _ := os.Hostname()

	switch exporter {
	case Jaeger:
		if len(endpoint) == 0 {
			return nil, fmt.Errorf("no endpoint defined for jaeger telemetry exporter")
		}
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint[0])))
	case StdOut:
		exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	// TODO: Figure out how to include CloudTrace without requiring necessary
	// configurations if not using it.
	// case CloudTrace:
	// 	exp, err = cloudTrace.New(cloudTrace.WithProjectID(projectID))
	default:
		return nil, fmt.Errorf("unknown exporter defined")
	}
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.TraceIDRatioBased(sampleRate/100)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
			attribute.String("hostname", hn),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}
