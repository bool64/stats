package stats

import "context"

// NoOp is a stats tracker stub.
type NoOp struct{}

var _ Tracker = NoOp{}

// Add discards value increment, can be negative.
func (NoOp) Add(ctx context.Context, name string, increment float64, labelsAndValues ...string) {}

// Set discards absolute value.
func (NoOp) Set(ctx context.Context, name string, absolute float64, labelsAndValues ...string) {}

// StatsTracker is a provider.
func (NoOp) StatsTracker() Tracker { return NoOp{} }
