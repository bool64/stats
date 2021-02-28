package stats_test

import (
	"context"
	"testing"

	"github.com/bool64/stats"
	"github.com/stretchr/testify/assert"
)

func TestAddKeysAndValues(t *testing.T) {
	ctx := stats.AddKeysAndValues(context.Background(), "k1", "one", "k2", "two")
	assert.Equal(t, []string{"k1", "one", "k2", "two"}, stats.KeysAndValues(ctx))

	ctx2 := stats.AddKeysAndValues(ctx, "k3", "three")
	assert.Equal(t, []string{"k1", "one", "k2", "two", "k3", "three"}, stats.KeysAndValues(ctx2))

	ctx3 := stats.AddKeysAndValues(ctx2, "k4", "four")
	assert.Equal(t, []string{"k1", "one", "k2", "two", "k3", "three", "k4", "four"}, stats.KeysAndValues(ctx3))

	ctx4 := stats.AddKeysAndValues(ctx2, "k4a", "four-a")
	assert.Equal(t, []string{"k1", "one", "k2", "two", "k3", "three", "k4a", "four-a"}, stats.KeysAndValues(ctx4))

	assert.Equal(t, []string{"k1", "one", "k2", "two", "k3", "three", "k4", "four"}, stats.KeysAndValues(ctx3))
}

func BenchmarkAddKeysAndValues(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := stats.AddKeysAndValues(context.Background(), "k1", "one", "k2", "two")
		_ = stats.KeysAndValues(ctx)

		ctx2 := stats.AddKeysAndValues(ctx, "k3", "three")
		_ = stats.KeysAndValues(ctx2)

		ctx3 := stats.AddKeysAndValues(ctx2, "k4", "four")
		_ = stats.KeysAndValues(ctx3)

		ctx4 := stats.AddKeysAndValues(ctx2, "k4a", "four-a")
		_ = stats.KeysAndValues(ctx4)
	}
}
