package stats

import "context"

// Tracker defines stats collector.
type Tracker interface {
	Adder
	Setter
}

// Adder defines incremental metric collector.
type Adder interface {
	// Add collects additional or observable value.
	Add(ctx context.Context, name string, increment float64, labelsAndValues ...string)
}

// AdderFunc implements Adder.
type AdderFunc func(ctx context.Context, name string, increment float64, labelsAndValues ...string)

// Add collects additional or observable value.
func (f AdderFunc) Add(ctx context.Context, name string, increment float64, labelsAndValues ...string) {
	f(ctx, name, increment, labelsAndValues...)
}

// Setter defines absolute value collector.
type Setter interface {
	// Set collects absolute value, e.g. number of goroutines.
	Set(ctx context.Context, name string, absolute float64, labelsAndValues ...string)
}

// SetterFunc implements Setter.
type SetterFunc func(ctx context.Context, name string, absolute float64, labelsAndValues ...string)

// Set collects absolute value.
func (f SetterFunc) Set(ctx context.Context, name string, absolute float64, labelsAndValues ...string) {
	f(ctx, name, absolute, labelsAndValues...)
}
