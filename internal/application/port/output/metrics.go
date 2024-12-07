package output

type MetricsReporter interface {
    IncrementCounter(name string, tags ...string)
    Gauge(name string, value float64, tags ...string)
    StartTimer(name string, tags ...string) Timer
}

type Timer interface {
    Stop()
    Duration() float64
} 