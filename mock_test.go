package stats_test

import (
	"context"
	"testing"

	"github.com/bool64/stats"
	"github.com/stretchr/testify/assert"
)

func TestTrackerMock(t *testing.T) {
	st := &stats.TrackerMock{}
	m1 := "metric1"
	m2 := "metric2"
	s1 := "strate1"

	st.Add(context.Background(), m1, 22, "not-relevant1", "abc1")
	st.Add(context.Background(), m1, 33, "not-relevant2", "abc2")
	st.Add(context.Background(), m2, 11, "not-relevant3", "abc3")
	st.Set(context.Background(), s1, 11, "not-relevant4", "abc4")
	st.Set(context.Background(), s1, 12, "not-relevant5", "abc5")

	assert.Equal(t, 55, st.Int(m1))
	assert.Equal(t, 11, st.Int(m2))
	assert.Equal(t, 12, st.Int(s1))

	assert.Equal(t, 55.0, st.Value(m1))
	assert.Equal(t, 11.0, st.Value(m2))
	assert.Equal(t, 12.0, st.Value(s1))

	assert.Equal(t, map[string]float64{m1: 55, m2: 11, s1: 12}, st.Values())

	exp := map[string]float64{
		"metric1{not-relevant1=\"abc1\"}": 22,
		"metric1{not-relevant2=\"abc2\"}": 33,
		"metric2{not-relevant3=\"abc3\"}": 11,
		"strate1{not-relevant4=\"abc4\"}": 11,
		"strate1{not-relevant5=\"abc5\"}": 12,
	}
	assert.Equal(t, exp, st.LabeledValues())

	assert.Equal(t, `metric1{not-relevant1="abc1"} 22
metric1{not-relevant2="abc2"} 33
metric2{not-relevant3="abc3"} 11
strate1{not-relevant4="abc4"} 11
strate1{not-relevant5="abc5"} 12`, st.Metrics(), st.Metrics())
}

func TestTrackerMock_nilValues(t *testing.T) {
	st := &stats.TrackerMock{}
	st.Add(context.Background(), "any", 22, "not-relevant1", "abc1")
	assert.Equal(t, 22, st.Int("any"))

	st = &stats.TrackerMock{}
	st.Set(context.Background(), "any", 22, "not-relevant1", "abc1")
	assert.Equal(t, 22, st.Int("any"))

	st = &stats.TrackerMock{}
	assert.Equal(t, 0, st.Int("none"))

	st = &stats.TrackerMock{}
	assert.Equal(t, 0.0, st.Value("none"))

	st = &stats.TrackerMock{}
	assert.Equal(t, map[string]float64{}, st.Values())
}

func TestTrackerMock_StatsTracker(t *testing.T) {
	st := &stats.TrackerMock{}
	assert.Equal(t, st, st.StatsTracker())
}
