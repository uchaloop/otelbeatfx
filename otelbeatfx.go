// Package otelbeatfx provides an OpenTelemetry beat.Handler through Uber Fx.
// The application owns the metric.MeterProvider, exporter and provider shutdown.
package otelbeatfx

import (
	"github.com/uchaloop/beat"
	"github.com/uchaloop/otelbeat"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// Module provides a beat.Handler using the container's metric.MeterProvider.
// It does not run work or configure the SDK, Resource, exporter or shutdown.
func Module() fx.Option {
	return fx.Module("otelbeatfx", fx.Provide(
		func(provider metric.MeterProvider) (beat.Handler, error) {
			return otelbeat.MakeHandler(provider.Meter("github.com/uchaloop/otelbeat"))
		},
	))
}
