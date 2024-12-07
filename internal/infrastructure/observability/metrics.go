package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsReporter struct {
	counters   map[string]prometheus.Counter
	histograms map[string]prometheus.Histogram
	registry   *prometheus.Registry
}

func NewMetricsReporter() *MetricsReporter {
	registry := prometheus.NewRegistry()
	
	return &MetricsReporter{
		counters:   make(map[string]prometheus.Counter),
		histograms: make(map[string]prometheus.Histogram),
		registry:   registry,
	}
}

func (r *MetricsReporter) RegisterCounter(name, help string, labels ...string) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	r.registry.MustRegister(counter)
	r.counters[name] = counter
}

func (r *MetricsReporter) RegisterHistogram(name, help string, buckets []float64, labels ...string) {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	})
	r.registry.MustRegister(histogram)
	r.histograms[name] = histogram
}

func (r *MetricsReporter) IncrementCounter(name string) {
	if counter, ok := r.counters[name]; ok {
		counter.Inc()
	}
}

type Timer struct {
	histogram prometheus.Histogram
	start    time.Time
}

func (r *MetricsReporter) StartTimer(name string) *Timer {
	if histogram, ok := r.histograms[name]; ok {
		return &Timer{
			histogram: histogram,
			start:    time.Now(),
		}
	}
	return nil
}

func (t *Timer) Stop() {
	if t != nil {
		duration := time.Since(t.start).Seconds()
		t.histogram.Observe(duration)
	}
} 