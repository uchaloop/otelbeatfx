package otelbeatfx_test

import (
	"context"
	"fmt"

	"github.com/uchaloop/beat"
	"github.com/uchaloop/otelbeatfx"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

func ExampleModule() {
	var observer beat.Handler
	app := fx.New(
		fx.Provide(func(lc fx.Lifecycle) (metric.MeterProvider, error) {
			// Read OTEL_SERVICE_NAME and OTEL_RESOURCE_ATTRIBUTES from the environment.
			res, err := resource.New(context.Background(), resource.WithFromEnv())
			if err != nil {
				return nil, err
			}
			// Configure an exporter-backed reader in the application.
			provider := sdkmetric.NewMeterProvider(sdkmetric.WithResource(res))
			lc.Append(fx.Hook{OnStop: provider.Shutdown})
			return provider, nil
		}),
		otelbeatfx.Module(),
		fx.Populate(&observer),
		fx.NopLogger,
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("handler available:", observer != nil)
	if err := app.Stop(context.Background()); err != nil {
		panic(err)
	}
	// Output: handler available: true
}
