package prom_test

import (
	"context"

	"github.com/bool64/stats"
	"github.com/bool64/stats/prom"
	"github.com/prometheus/client_golang/prometheus"
)

func ExampleTracker() {
	// Bring your own Prometheus registry.
	registry := prometheus.NewRegistry()
	tr := prom.Tracker{
		Registry: registry,
	}

	// Add custom Prometheus configuration where necessary.
	tr.DeclareHistogram("my_latency_seconds", prometheus.HistogramOpts{
		Buckets: []float64{1e-4, 1e-3, 1e-2, 1e-1, 1, 10, 100},
	})

	ctx := context.Background()

	// Add labels to context.
	ctx = stats.AddKeysAndValues(ctx, "ctx-label", "ctx-value0")

	// Override label values.
	ctx = stats.AddKeysAndValues(ctx, "ctx-label", "ctx-value1")

	// Collect stats with last mile labels.
	tr.Add(ctx, "my_count", 1,
		"some-label", "some-value",
	)

	tr.Add(ctx, "my_latency_seconds", 1.23)

	tr.Set(ctx, "temperature", 33.3)
}
