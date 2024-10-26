package meter

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

var observers = []func(meter metric.Meter) error{
	observeGaugeMemoryUsage,
	observeCounterNumGoRoutines,
}

var memoryUsage = func() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Alloc)
}

var observeGaugeMemoryUsage = func(meter metric.Meter) error {
	_, err := meter.Int64ObservableGauge("memory_usage_bytes",
		metric.WithDescription("Memory Usage in Bytes"),
		metric.WithUnit("bytes"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(memoryUsage())
			return nil
		}))
	if err != nil {
		return err
	}

	return nil
}

var observeCounterNumGoRoutines = func(meter metric.Meter) error {
	_, err := meter.Int64ObservableCounter("num_go_routines",
		metric.WithDescription("Number of Go Routines"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(runtime.NumGoroutine()))
			return nil
		}))
	if err != nil {
		return err
	}

	return nil
}
