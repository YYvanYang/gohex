package config

import (
	"time"
	"github.com/spf13/viper"
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Driver          string        `yaml:"driver"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime
}

"`func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

type AuthConfig struct {
	JWT struct {
		SecretKey      string        `yaml:"secret_key"`
		AccessTTL      time.Duration `yaml:"access_ttl"`
		RefreshTTL     time.Duration `yaml:"refresh_ttl"`
		Issuer        string        `yaml:"issuer"`
		Audience      []string      `yaml:"audience"`
		SigningMethod string        `yaml:"signing_method"`
	} `yaml:"jwt"`

	Password struct {
		MinLength      int    `yaml:"min_length"`
		HashAlgorithm  string `yaml:"hash_algorithm"`
		HashMemory     int    `yaml:"hash_memory"`
		HashIterations int    `yaml:"hash_iterations"`
	} `yaml:"password"`

	Session struct {
		Enabled      bool          `yaml:"enabled"`
		Store        string        `yaml:"store"`
		CookieName   string        `yaml:"cookie_name"`
		CookiePath   string        `yaml:"cookie_path"`
		CookieDomain string        `yaml:"cookie_domain"`
		MaxAge       time.Duration `yaml:"max_age"`
		Secure       bool          `yaml:"secure"`
		HttpOnly     bool          `yaml:"http_only"`
	} `yaml:"session"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	viper.SetConfigFile(configPath)
	viper.SetEnvPrefix("GOHEX")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.App.Name == "" {
		return errors.New("app name is required")
	}
	if c.App.Environment == "" {
		return errors.New("app environment is required")
	}
	if c.HTTP.Port <= 0 {
		return errors.New("invalid HTTP port")
	}
	if c.HTTP.ReadTimeout <= 0 {
		return errors.New("invalid HTTP read timeout")
	}
	if c.HTTP.WriteTimeout <= 0 {
		return errors.New("invalid HTTP write timeout")
	}
	if c.Database.Driver == "" {
		return errors.New("database driver is required")
	}
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	if c.Database.Port <= 0 {
		return errors.New("invalid database port")
	}
	if c.Database.MaxOpenConns <= 0 {
		return errors.New("invalid max open connections")
	}
	if c.JWT.SecretKey == "" {
		return errors.New("JWT secret key is required")
	}
	if c.JWT.TokenDuration <= 0 {
		return errors.New("invalid JWT token duration")
	}
	return nil
} 