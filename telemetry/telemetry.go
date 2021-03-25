package telemetry

import "github.com/corioders/gokit/application"

type SetupTelemetryOptions struct {
	ServiceName string

	TraceExporter *TraceExporterOptions
}

func (o *SetupTelemetryOptions) validate() error {
	return o.TraceExporter.validate()
}

func StartTelemetry(app application.StopHandler, options *SetupTelemetryOptions) error {
	err := options.validate()
	if err != nil {
		return err
	}

	switch options.TraceExporter.ExporeterType {
	case TraceExpoterJaeger:
		flush, err := StartJaegerTracing(options.ServiceName, options.TraceExporter.CollectorEndpoint)
		app.StopFunc("Flush jaeger tracing exporter", func() error {
			flush()
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
