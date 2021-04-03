package telemetry

import "github.com/corioders/gokit/application"

type SetupTelemetryOptions struct {
	ServiceName string

	TraceExporter *TraceExporterOptions
}

func (o *SetupTelemetryOptions) validate() error {
	return o.TraceExporter.validate()
}

func StartTelemetry(sr application.StopRegistrar, options *SetupTelemetryOptions) error {
	err := options.validate()
	if err != nil {
		return err
	}

	switch options.TraceExporter.ExporeterType {
	case TraceExpoterJaeger:
		flush, err := StartJaegerTracing(options.ServiceName, options.TraceExporter.CollectorEndpoint)
		sr.RegisterOnStop("Flush jaeger tracing exporter", func() error {
			flush()
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
