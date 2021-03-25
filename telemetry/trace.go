package telemetry

import (
	"errors"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

type TraceExporterOptions struct {
	CollectorEndpoint string
	ExporeterType     traceExporterType
}

func (o *TraceExporterOptions) validate() error {
	if o.ExporeterType == 0 {
		return errors.New("Trace exporter type not provided")
	}
	if o.CollectorEndpoint == "" {
		return errors.New("Trace collector endpoint not provided")
	}
	return nil
}

type traceExporterType int

const (
	TraceExpoterJaeger traceExporterType = iota + 1
)

// StartJaegerTracing starts opentelemety jaeger exporter and sets it as global tracing exporter.
func StartJaegerTracing(serviceName, collectorEndpoint string) (flush func(), err error) {
	serviceNameResource := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.ServiceNameKey.String(serviceName)))

	// Create and install Jaeger export pipeline.
	return jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(collectorEndpoint),
		jaeger.WithSDKOptions(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithResource(serviceNameResource),
		),
	)
}
