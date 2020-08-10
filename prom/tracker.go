package prom

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/bool64/stats"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	_ stats.Tracker         = &Tracker{}
	_ stats.TrackerProvider = &Tracker{}
)

// NewStatsTracker creates prometheus stats tracker.
func NewStatsTracker(registry *prometheus.Registry) (*Tracker, error) {
	if registry == nil {
		registry = prometheus.NewRegistry()

		err := registry.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		if err != nil {
			return nil, err
		}

		err = registry.Register(prometheus.NewGoCollector())
		if err != nil {
			return nil, err
		}
	}

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	t := &Tracker{
		Registry: registry,

		reg: reg,
	}

	return t, nil
}

// Tracker implements stats tracker with prometheus registry.
//
// Please use NewStatsTracker to create new instance.
type Tracker struct {
	mu         sync.Mutex
	collectors map[identity]func(map[string]string, float64)
	histograms map[string]prometheus.HistogramOpts
	summaries  map[string]prometheus.SummaryOpts
	counters   map[string]prometheus.CounterOpts
	gauges     map[string]prometheus.GaugeOpts

	reg *regexp.Regexp

	Registry  *prometheus.Registry
	ErrLogger func(ctx context.Context, err error, labels []string)
}

// StatsTracker is a service locator provider.
func (t *Tracker) StatsTracker() stats.Tracker {
	return t
}

// Add collects value to Counter, Summary or Histogram.
func (t *Tracker) Add(ctx context.Context, name string, value float64, labelsAndValues ...string) {
	t.collect(ctx, false, name, value, labelsAndValues)
}

// Set collects absolute value to Gauge.
func (t *Tracker) Set(ctx context.Context, name string, absolute float64, labelsAndValues ...string) {
	t.collect(ctx, true, name, absolute, labelsAndValues)
}

// DeclareHistogram registers histogram metric for given name.
func (t *Tracker) DeclareHistogram(name string, opts prometheus.HistogramOpts) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.histograms == nil {
		t.histograms = make(map[string]prometheus.HistogramOpts)
	}

	t.histograms[t.prepareName(name)] = opts
}

// DeclareSummary registers summary metric for given name.
func (t *Tracker) DeclareSummary(name string, opts prometheus.SummaryOpts) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.summaries == nil {
		t.summaries = make(map[string]prometheus.SummaryOpts)
	}

	t.summaries[t.prepareName(name)] = opts
}

// DeclareCounter registers counter metric for given name.
func (t *Tracker) DeclareCounter(name string, opts prometheus.CounterOpts) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.counters == nil {
		t.counters = make(map[string]prometheus.CounterOpts)
	}

	t.counters[t.prepareName(name)] = opts
}

// DeclareGauge registers gauge metric for given name.
func (t *Tracker) DeclareGauge(name string, opts prometheus.GaugeOpts) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.gauges == nil {
		t.gauges = make(map[string]prometheus.GaugeOpts)
	}

	t.gauges[t.prepareName(name)] = opts
}

type identity struct {
	name      string
	xorLabels string
}

var labelBufferPool = sync.Pool{
	New: func() interface{} {
		return &labelBuffer{
			values:    make(map[string]string),
			xorLabels: make([]byte, 0),
		}
	},
}

type labelBuffer struct {
	labelsAndValues []string
	values          map[string]string
	xorLabels       []byte
}

func (lb *labelBuffer) reset() {
	lb.labelsAndValues = lb.labelsAndValues[:0]
	lb.xorLabels = lb.xorLabels[:0]

	for k := range lb.values {
		delete(lb.values, k)
	}
}

func (t *Tracker) collect(ctx context.Context, isGauge bool, name string, value float64, labelsAndValues []string) {
	lb, ok := labelBufferPool.Get().(*labelBuffer)
	if !ok {
		panic("BUG: could not assert *labelBuffer type")
	}

	lb.reset()
	defer labelBufferPool.Put(lb)

	toLabels(ctx, lb, labelsAndValues)

	id := identity{
		name:      name,
		xorLabels: byteSlice2String(lb.xorLabels),
	}

	t.mu.Lock()

	collect, ok := t.collectors[id]
	if !ok {
		// If new entry to collectors is to be created, immutable string is needed.
		// Otherwise if unsafe string is kept, map keys will change along with sync.Pool reused []byte values.
		id.xorLabels = string(lb.xorLabels)

		labels := make([]string, 0, len(lb.values))
		for l := range lb.values {
			labels = append(labels, l)
		}

		sort.Strings(labels)

		// Canonical identity contains prepared name.
		canonicalID := identity{
			name:      t.prepareName(name),
			xorLabels: id.xorLabels,
		}

		collect, ok = t.collectors[canonicalID]
		_ = ok ||
			t.histogram(ctx, canonicalID, labels, &collect) ||
			t.summary(ctx, canonicalID, labels, &collect) ||
			(!isGauge && t.counter(ctx, canonicalID, labels, &collect)) ||
			t.gauge(ctx, canonicalID, labels, &collect)

		if t.collectors == nil {
			t.collectors = make(map[identity]func(map[string]string, float64))
		}

		t.collectors[id] = collect

		if id != canonicalID {
			t.collectors[canonicalID] = collect
		}
	}
	t.mu.Unlock()

	collect(lb.values, value)
}

func (t *Tracker) summary(
	ctx context.Context,
	canonicalID identity,
	labels []string,
	collect *func(map[string]string, float64),
) bool {
	opts, ok := t.summaries[canonicalID.name]
	if !ok {
		return false
	}

	opts.Name = coalesce(opts.Name, canonicalID.name)
	opts.Help = coalesce(opts.Help, "is created by "+callerFunc(4))

	summary := prometheus.NewSummaryVec(opts, labels)

	err := t.Registry.Register(summary)
	if err != nil && t.ErrLogger != nil {
		t.ErrLogger(ctx, err, labels)
	}

	*collect = func(labelValues map[string]string, value float64) {
		summary.With(labelValues).Observe(value)
	}

	return true
}

func (t *Tracker) histogram(
	ctx context.Context,
	canonicalID identity,
	labels []string,
	collect *func(map[string]string, float64),
) bool {
	opts, ok := t.histograms[canonicalID.name]
	if !ok {
		return false
	}

	opts.Name = coalesce(opts.Name, canonicalID.name)
	opts.Help = coalesce(opts.Help, "is created by "+callerFunc(4))

	histogram := prometheus.NewHistogramVec(opts, labels)

	err := t.Registry.Register(histogram)
	if err != nil && t.ErrLogger != nil {
		t.ErrLogger(ctx, err, labels)
	}

	*collect = func(labelValues map[string]string, value float64) {
		histogram.With(labelValues).Observe(value)
	}

	return true
}

func (t *Tracker) counter(
	ctx context.Context,
	canonicalID identity,
	labels []string,
	collect *func(map[string]string, float64),
) bool {
	opts := t.counters[canonicalID.name]

	opts.Name = coalesce(opts.Name, canonicalID.name)
	opts.Help = coalesce(opts.Help, "is created by "+callerFunc(4))

	counter := prometheus.NewCounterVec(opts, labels)

	err := t.Registry.Register(counter)
	if err != nil && t.ErrLogger != nil {
		t.ErrLogger(ctx, err, labels)
	}

	*collect = func(labelValues map[string]string, value float64) {
		counter.With(labelValues).Add(value)
	}

	return true
}

func (t *Tracker) gauge(
	ctx context.Context,
	canonicalID identity,
	labels []string,
	collect *func(map[string]string, float64),
) bool {
	opts := t.gauges[canonicalID.name]

	opts.Name = coalesce(opts.Name, canonicalID.name)
	opts.Help = coalesce(opts.Help, "is created by "+callerFunc(4))

	gauge := prometheus.NewGaugeVec(opts, labels)

	err := t.Registry.Register(gauge)
	if err != nil && t.ErrLogger != nil {
		t.ErrLogger(ctx, err, labels)
	}

	*collect = func(labelValues map[string]string, value float64) {
		gauge.With(labelValues).Set(value)
	}

	return true
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}

	return b
}

func toLabels(ctx context.Context, lb *labelBuffer, labelsAndValues []string) {
	ctxKV := stats.KeysAndValues(ctx)
	if len(ctxKV) > 0 {
		lb.labelsAndValues = append(lb.labelsAndValues, ctxKV...)
	}

	lb.labelsAndValues = append(lb.labelsAndValues, labelsAndValues...)
	label := ""

	for _, l := range lb.labelsAndValues {
		if label == "" {
			label = l
		} else {
			for i, c := range []byte(label) {
				if i >= len(lb.xorLabels) {
					lb.xorLabels = append(lb.xorLabels, c)
				} else {
					lb.xorLabels[i] ^= c
				}
			}

			lb.values[label] = l
			label = ""
		}
	}
}

func (t *Tracker) prepareName(name string) string {
	return strings.Trim(t.reg.ReplaceAllString(name, "_"), "_")
}
