package stats

import (
	"context"
	"sync"
)

// TrackerMock can collect stats for tests with labels ignored.
type TrackerMock struct {
	mu     sync.Mutex
	values map[string]float64
}

// Add collects metric increment.
func (t *TrackerMock) Add(_ context.Context, name string, increment float64, _ ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		t.values = make(map[string]float64, 10)
	}

	t.values[name] += increment
}

// Set collects absolute value.
func (t *TrackerMock) Set(_ context.Context, name string, absolute float64, _ ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		t.values = make(map[string]float64, 10)
	}

	t.values[name] = absolute
}

// Value returns collected value by name.
func (t *TrackerMock) Value(name string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		return 0
	}

	return t.values[name]
}

// Int returns collected value as integer by name.
func (t *TrackerMock) Int(name string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		return 0
	}

	return int(t.values[name])
}

// Values returns collected values as a map.
func (t *TrackerMock) Values() map[string]float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		return map[string]float64{}
	}

	res := make(map[string]float64, len(t.values))

	for k, v := range t.values {
		res[k] = v
	}

	return res
}

// StatsTracker is a provider.
func (t *TrackerMock) StatsTracker() Tracker {
	return t
}
