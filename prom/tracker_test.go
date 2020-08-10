package prom_test

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/bool64/stats"
	"github.com/bool64/stats/prom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

const expected = `
		# HELP another_size is another size
		# TYPE another_size summary
		another_size_sum{l="v"} 4950
		another_size_count{l="v"} 100
		# HELP iterations is created by stats/prom_test.TestTracker.func1
		# TYPE iterations counter
		iterations 100
		# HELP short is short
		# TYPE short gauge
		short 1234
		# HELP some_action_count is counting some actions
		# TYPE some_action_count counter
		some_action_count{l="v",name="name-val",other="another-val"} 1300
		some_action_count{l="v",name="name-val",other="other-val"} 2500
		# HELP some_action_items is created by stats/prom_test.TestTracker.func1
		# TYPE some_action_items gauge
		some_action_items{l="v",name="name-val",other="other-val"} 123
		# HELP some_size is some size
		# TYPE some_size histogram
		some_size_bucket{l="v",le="1"} 2
		some_size_bucket{l="v",le="5"} 6
		some_size_bucket{l="v",le="10"} 11
		some_size_bucket{l="v",le="50"} 51
		some_size_bucket{l="v",le="100"} 100
		some_size_bucket{l="v",le="+Inf"} 100
		some_size_sum{l="v"} 4950
		some_size_count{l="v"} 100
	`

func TestTracker(t *testing.T) {
	registry := prometheus.NewPedanticRegistry()
	tr, err := prom.NewStatsTracker(registry)
	assert.NoError(t, err)

	tr.DeclareHistogram("some.size", prometheus.HistogramOpts{
		Help:    "is some size",
		Buckets: []float64{1, 5, 10, 50, 100},
	})
	tr.DeclareSummary("another.size", prometheus.SummaryOpts{Help: "is another size"})
	tr.DeclareCounter("some_action_count", prometheus.CounterOpts{Help: "is counting some actions"})
	tr.DeclareGauge("short", prometheus.GaugeOpts{Help: "is short"})

	wg := sync.WaitGroup{}

	ctx := context.Background()

	ctx = stats.AddKeysAndValues(ctx, "l", "0")
	ctx = stats.AddKeysAndValues(ctx, "l", "v")

	for i := 0; i < 100; i++ {
		i := i

		wg.Add(1)

		go func() {
			defer wg.Done()
			tr.Add(context.Background(), "iterations", 1)
			tr.Add(ctx, "some.action.count", 12,
				"name", "name-val",
				"other", "other-val",
			)
			tr.Add(ctx, "some.action.count", 13,
				"name", "name-val",
				"other", "other-val",
			)

			tr.Add(ctx, "some.action.count", 13,
				"name", "name-val",
				"other", "another-val",
			)
			tr.Set(ctx, "some.action.items", 123,
				"name", "name-val",
				"other", "other-val",
			)

			tr.Add(ctx, "some.size", float64(i))
			tr.Add(ctx, "another.size", float64(i))
		}()
	}

	tr.Set(context.Background(), "short", 1234)

	wg.Wait()

	err = testutil.GatherAndCompare(registry, bytes.NewBufferString(expected))
	assert.NoError(t, err)
}
