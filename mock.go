package stats

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// TrackerMock can collect stats for tests with labels ignored.
type TrackerMock struct {
	mu            sync.Mutex
	values        map[string]float64
	labeledValues map[string]float64
}

var escaper = strings.NewReplacer("\n", `\\n`, `\`, `\\`, `"`, `\"`)

func labelsString(labels []string) string {
	if len(labels) == 0 {
		return ""
	}

	isKey := true
	key := ""

	if len(labels)%2 != 0 {
		panic("malformed pairs")
	}

	type kv struct {
		k, v string
	}

	lb := make([]kv, 0, len(labels)/2)

	for _, l := range labels {
		if isKey {
			if l == "" {
				panic("empty key received in labels")
			}

			key = l
			isKey = false
		} else {
			lb = append(lb, kv{k: key, v: l})
			isKey = true
		}
	}

	sort.Slice(lb, func(i, j int) bool {
		return lb[i].k < lb[j].v
	})

	res := ""
	for _, i := range lb {
		res += i.k + `="` + escaper.Replace(i.v) + `",`
	}

	return res[0 : len(res)-1]
}

// Add collects metric increment.
func (t *TrackerMock) Add(_ context.Context, name string, increment float64, labels ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		t.values = make(map[string]float64)
		t.labeledValues = make(map[string]float64)
	}

	t.values[name] += increment
	t.labeledValues[name+"{"+labelsString(labels)+"}"] += increment
}

// Set collects absolute value.
func (t *TrackerMock) Set(_ context.Context, name string, absolute float64, labels ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		t.values = make(map[string]float64)
		t.labeledValues = make(map[string]float64)
	}

	t.values[name] = absolute
	t.labeledValues[name+"{"+labelsString(labels)+"}"] = absolute
}

// Value returns collected value by name.
func (t *TrackerMock) Value(name string, labels ...string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		return 0
	}

	if len(labels) > 0 {
		return t.labeledValues[name+"{"+labelsString(labels)+"}"]
	}

	return t.values[name]
}

// Int returns collected value as integer by name.
func (t *TrackerMock) Int(name string, labels ...string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.values == nil {
		return 0
	}

	if len(labels) > 0 {
		return int(t.labeledValues[name+"{"+labelsString(labels)+"}"])
	}

	return int(t.values[name])
}

// Values returns collected summarized values as a map.
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

// LabeledValues returns collected labeled values as a map.
func (t *TrackerMock) LabeledValues() map[string]float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.labeledValues == nil {
		return map[string]float64{}
	}

	res := make(map[string]float64, len(t.labeledValues))

	for k, v := range t.labeledValues {
		res[k] = v
	}

	return res
}

// Metrics returns collected values in Prometheus format.
func (t *TrackerMock) Metrics() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.labeledValues == nil {
		return ""
	}

	res := make([]string, 0, len(t.labeledValues))
	for k, v := range t.labeledValues {
		res = append(res, k+" "+strconv.FormatFloat(v, 'g', -1, 64))
	}

	sort.Strings(res)

	return strings.Join(res, "\n")
}

// StatsTracker is a provider.
func (t *TrackerMock) StatsTracker() Tracker {
	return t
}
