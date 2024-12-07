package query

type MiddlewareConfig interface {
    // 启用/禁用中间件
    EnableValidation() bool
    EnableCaching() bool
    EnableLogging() bool
    EnableMetrics() bool

    // 中间件配置
    ValidationConfig() ValidationConfig
    CacheConfig() CacheConfig
    LoggingConfig() LoggingConfig
    MetricsConfig() MetricsConfig
}

type ValidationConfig struct {
    SkipTypes []interface{}  // 跳过验证的类型
}

type CacheConfig struct {
    DefaultTTL time.Duration
    MaxSize    int
}

type LoggingConfig struct {
    LogLevel string
    Fields   map[string]interface{}
}

type MetricsConfig struct {
    Namespace string
    Subsystem string
} 