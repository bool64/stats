package stats_test

import (
	"context"
	"testing"

	"github.com/bool64/stats"
	"github.com/stretchr/testify/assert"
)

func TestNoOp_Metric(_ *testing.T) {
	stats.NoOp{}.Add(context.Background(), "any", 1, "key", "val")
}

func TestNoOp_State(_ *testing.T) {
	stats.NoOp{}.Set(context.Background(), "any", 1, "key", "val")
}

func TestNoOp_StatsTracker(t *testing.T) {
	_, ok := stats.NoOp{}.StatsTracker().(stats.NoOp)
	assert.True(t, ok)
}
