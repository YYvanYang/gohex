package metrics

import (
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gohex/gohex/internal/application/port"
)

type prometheusMetrics struct {
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
	timers   map[string]*prometheus.HistogramVec
}

func NewPrometheusMetrics(namespace string) MetricsReporter {
	return &prometheusMetrics{
		counters: make(map[string]*prometheus.CounterVec),
		gauges:   make(map[string]*prometheus.GaugeVec),
		timers:   make(map[string]*prometheus.HistogramVec),
	}
}

func (m *prometheusMetrics) IncrementCounter(name string, tags ...string) {
	counter, ok := m.counters[name]
	if !ok {
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
			},
			extractTagKeys(tags),
		)
		prometheus.MustRegister(counter)
		m.counters[name] = counter
	}
	counter.WithLabelValues(extractTagValues(tags)...).Inc()
}

func (m *prometheusMetrics) Gauge(name string, value float64, tags ...string) {
	gauge, ok := m.gauges[name]
	if !ok {
		gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
			},
			extractTagKeys(tags),
		)
		prometheus.MustRegister(gauge)
		m.gauges[name] = gauge
	}
	gauge.WithLabelValues(extractTagValues(tags)...).Set(value)
}

func (m *prometheusMetrics) StartTimer(name string, tags ...string) Timer {
	timer, ok := m.timers[name]
	if !ok {
		timer = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    name,
				Buckets: prometheus.DefBuckets,
			},
			extractTagKeys(tags),
		)
		prometheus.MustRegister(timer)
		m.timers[name] = timer
	}
	
	start := time.Now()
	return &prometheusTimer{
		start:     start,
		histogram: timer,
		tags:      tags,
	}
}

type prometheusTimer struct {
	start     time.Time
	histogram *prometheus.HistogramVec
	tags      []string
}

func (t *prometheusTimer) Stop() {
	duration := time.Since(t.start).Seconds()
	t.histogram.WithLabelValues(extractTagValues(t.tags)...).Observe(duration)
}

func (t *prometheusTimer) Duration() float64 {
	return time.Since(t.start).Seconds()
}

// 添加辅助方法
func extractTagKeys(tags []string) []string {
	keys := make([]string, len(tags)/2)
	for i := 0; i < len(tags); i += 2 {
		keys[i/2] = tags[i]
	}
	return keys
}

func extractTagValues(tags []string) []string {
	values := make([]string, len(tags)/2)
	for i := 1; i < len(tags); i += 2 {
		values[(i-1)/2] = tags[i]
	}
	return values
}

// 实现其他方法... 