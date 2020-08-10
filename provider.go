package stats

// TrackerProvider defines service locator interface.
type TrackerProvider interface {
	StatsTracker() Tracker
}
