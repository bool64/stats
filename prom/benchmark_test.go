// Sample result:
//   BenchmarkTracker_Add-16      	 1700233	       728 ns/op	       0 B/op	       0 allocs/op
//   BenchmarkRawPrometheus-16    	 1564921	       776 ns/op	     336 B/op	       2 allocs/op

package prom_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/bool64/stats"
	"github.com/bool64/stats/prom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func BenchmarkTracker_Add(b *testing.B) {
	registry := prometheus.NewRegistry()
	tr, err := prom.NewStatsTracker(registry)
	require.NoError(b, err)

	ctx := context.Background()
	ctx = stats.AddKeysAndValues(ctx, "foo", "bar")

	var lv []string
	for i := 0; i < 20; i++ {
		lv = append(lv, "some-val"+strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	sem := make(chan struct{}, 50)

	for i := 0; i < b.N; i++ {
		i := i

		sem <- struct{}{}

		go func() {
			defer func() {
				<-sem
			}()

			tr.Add(ctx, "some.action.count", 1,
				"name", lv[i%10],
				"other", lv[i%20],
			)
		}()
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}

func BenchmarkRawPrometheus(b *testing.B) {
	registry := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "namespace_some_action_count",
		},
		[]string{"name", "other"},
	)
	registry.MustRegister(counter)

	sem := make(chan struct{}, 50)

	var lv []string
	for i := 0; i < 20; i++ {
		lv = append(lv, "some-val"+strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		i := i

		sem <- struct{}{}

		go func() {
			defer func() {
				<-sem
			}()

			counter.With(prometheus.Labels{
				"name":  lv[i%10],
				"other": lv[i%20],
			}).Add(1)
		}()
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}
