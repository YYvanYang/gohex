package config

type CommandBusConfig struct {
    Middleware CommandMiddlewareConfig `yaml:"middleware"`
    Handlers   CommandHandlerConfig   `yaml:"handlers"`
    Metrics    CommandMetricsConfig   `yaml:"metrics"`
}

type CommandHandlerConfig struct {
    // 处理器超时配置
    Timeout time.Duration `yaml:"timeout"`
    // 处理器并发限制
    MaxConcurrency int `yaml:"max_concurrency"`
    // 处理器重试配置
    RetryAttempts int           `yaml:"retry_attempts"`
    RetryDelay    time.Duration `yaml:"retry_delay"`
}

type CommandMetricsConfig struct {
    // 指标前缀
    Namespace string `yaml:"namespace"`
    // 指标标签
    Labels map[string]string `yaml:"labels"`
    // 指标收集配置
    CollectInterval time.Duration `yaml:"collect_interval"`
    // 指标导出配置
    ExportPath string `yaml:"export_path"`
} 