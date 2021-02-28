# Contextualized Stats Tracker for Go

This library provides context-driven stats tracker.

[![Build Status](https://github.com/bool64/stats/workflows/test/badge.svg)](https://github.com/bool64/stats/actions?query=branch%3Amaster+workflow%3Atest)
[![Coverage Status](https://codecov.io/gh/bool64/stats/branch/master/graph/badge.svg)](https://codecov.io/gh/bool64/stats)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/bool64/stats)
[![time tracker](https://wakatime.com/badge/github/bool64/stats.svg)](https://wakatime.com/badge/github/bool64/stats)
![Code lines](https://sloc.xyz/github/bool64/stats/?category=code)
![Comments](https://sloc.xyz/github/bool64/stats/?category=comments)

## Features

* Loosely coupled with underlying implementation.
* Context-driven labels control.
* Zero allocation implementation for [Prometheus client](https://github.com/bool64/prom-stats).
* A simple interface with variadic number of key-value pairs for labels.
* Easily mockable interface free from 3rd party dependencies.

## Example

```go
// Bring your own Prometheus registry.
registry := prometheus.NewRegistry()
tr := prom.Tracker{
    Registry: registry,
}

// Add custom Prometheus configuration where necessary.
tr.DeclareHistogram("my_latency_seconds", prometheus.HistogramOpts{
    Buckets: []float64{1e-4, 1e-3, 1e-2, 1e-1, 1, 10, 100},
})

ctx := context.Background()

// Add labels to context.
ctx = stats.AddKeysAndValues(ctx, "ctx-label", "ctx-value0")

// Override label values.
ctx = stats.AddKeysAndValues(ctx, "ctx-label", "ctx-value1")

// Collect stats with last mile labels.
tr.Add(ctx, "my_count", 1,
    "some-label", "some-value",
)

tr.Add(ctx, "my_latency_seconds", 1.23)

tr.Set(ctx, "temperature", 33.3)
```

## Versioning

This project adheres to [Semantic Versioning](https://semver.org/#semantic-versioning-200).

Before version `1.0.0`, breaking changes are tagged with `MINOR` bump, features and fixes are tagged with `PATCH` bump.
After version `1.0.0`, breaking changes are tagged with `MAJOR` bump.
