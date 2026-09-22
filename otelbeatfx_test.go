package otelbeatfx_test

import (
	"context"
	"github.com/uchaloop/beat"
	"github.com/uchaloop/job"
	"github.com/uchaloop/otelbeatfx"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/fx"
	"testing"
	"time"
)

func TestModuleRecordsMetrics(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	var handler beat.Handler
	app := fx.New(
		fx.Provide(func() metric.MeterProvider { return provider }),
		otelbeatfx.Module(), fx.Populate(&handler), fx.NopLogger,
	)
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
	result := job.Result{Start: time.Now(), Duration: time.Second, WorkDuration: time.Second, Processed: 42, Outcome: job.OutcomeOK}
	handler.Handle(context.Background(), beat.Record{Result: result, Mode: beat.ModeFixedRate, ScheduledFor: result.Start})
	var data metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &data); err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, scope := range data.ScopeMetrics {
		if scope.Scope.Name != "github.com/uchaloop/otelbeat" {
			t.Fatalf("scope: %s", scope.Scope.Name)
		}
		for _, m := range scope.Metrics {
			if m.Name == "job.run.processed" {
				h := m.Data.(metricdata.Histogram[int64])
				if len(h.DataPoints) != 1 || h.DataPoints[0].Count != 1 || h.DataPoints[0].Sum != 42 {
					t.Fatalf("processed: %+v", h)
				}
				found = true
			}
		}
	}
	if !found {
		t.Fatal("processed measurement missing")
	}
}

func TestModuleRequiresProvider(t *testing.T) {
	var handler beat.Handler
	app := fx.New(otelbeatfx.Module(), fx.Populate(&handler), fx.NopLogger)
	if app.Err() == nil {
		t.Fatal("expected missing MeterProvider error")
	}
}
