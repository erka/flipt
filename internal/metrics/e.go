package metrics

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// exporter is an OpenTelemetry metric exporter.
type exporter struct {
	shutdownOnce sync.Once
}

// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New() (metric.Exporter, error) {
	exp := &exporter{}
	return exp, nil
}

func (e *exporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func (e *exporter) Aggregation(k metric.InstrumentKind) metric.Aggregation {
	return nil
}

func (e *exporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) error {
	select {
	case <-ctx.Done():
		// Don't do anything if the context has already timed out.
		return ctx.Err()
	default:
		// Context is still valid, continue.
		for _, sm := range data.ScopeMetrics {
			if sm.Scope.Name != "github.com/flipt-io/flipt" {
				continue
			}
			for _, m := range sm.Metrics {
				if m.Name != "flipt_evaluations_results" {
					continue
				}
				switch d := m.Data.(type) {
				case metricdata.Sum[int64]:
					fmt.Printf("%+v", d.DataPoints)
				}

			}
		}
	}
	return nil
}

func (e *exporter) ForceFlush(ctx context.Context) error {
	// exporter holds no state, nothing to flush.
	fmt.Println("ForceFlush")
	return ctx.Err()
}

func (e *exporter) Shutdown(ctx context.Context) error {
	e.shutdownOnce.Do(func() {

	})
	return ctx.Err()
}

func (e *exporter) MarshalLog() interface{} {
	return struct{ Type string }{Type: "STDOUT"}
}

// Collect implements prometheus.Collector.
//
// This method is safe to call concurrently.
func (e *exporter) Collect(ctx context.Context, rm *metricdata.ResourceMetrics) error {

	return nil
}
