package config

type CommandMiddlewareConfig struct {
    Validation   ValidationConfig   `yaml:"validation"`
    Transaction  TransactionConfig `yaml:"transaction"`
    Logging      LoggingConfig     `yaml:"logging"`
    Metrics      MetricsConfig     `yaml:"metrics"`
    Events       EventConfig       `yaml:"events"`
}

type TransactionConfig struct {
    Enabled     bool   `yaml:"enabled"`
    Propagation string `yaml:"propagation"` // Required, RequiresNew, Supports
    Timeout     time.Duration `yaml:"timeout"`
    Isolation   string `yaml:"isolation"`   // ReadCommitted, RepeatableRead, Serializable
}

type EventConfig struct {
    Enabled         bool `yaml:"enabled"`
    AsyncPublishing bool `yaml:"async_publishing"`
    BatchSize       int  `yaml:"batch_size"`
    RetryAttempts   int  `yaml:"retry_attempts"`
} 