package config

type QueryMiddlewareConfig struct {
    Validation ValidationConfig `yaml:"validation"`
    Cache      CacheConfig     `yaml:"cache"`
    Logging    LoggingConfig   `yaml:"logging"`
    Retry      RetryConfig     `yaml:"retry"`
    Timeout    TimeoutConfig   `yaml:"timeout"`
}

type RetryConfig struct {
    Enabled    bool          `yaml:"enabled"`
    MaxRetries int           `yaml:"max_retries"`
    Backoff    time.Duration `yaml:"backoff"`
}

type TimeoutConfig struct {
    Enabled bool          `yaml:"enabled"`
    Default time.Duration `yaml:"default"`
} 