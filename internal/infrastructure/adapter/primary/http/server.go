package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	router  *Router
	server  *http.Server
	logger  Logger
	metrics MetricsReporter
	db      Database
	cache   Cache
	mq      MessageQueue
	config  ServerConfig
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	Health       HealthCheckConfig `yaml:"health"`
}

type HealthCheckConfig struct {
	Enabled       bool          `yaml:"enabled"`
	Interval      time.Duration `yaml:"interval"`
	Timeout       time.Duration `yaml:"timeout"`
	InitialDelay  time.Duration `yaml:"initial_delay"`
	FailureThreshold int        `yaml:"failure_threshold"`
}

type healthChecker struct {
	config   HealthCheckConfig
	logger   Logger
	metrics  MetricsReporter
	checks   []Check
	failures int
	mu       sync.Mutex
}

type Check struct {
	Name     string
	Check    func(context.Context) error
	Required bool
}

func NewServer(
	cfg ServerConfig,
	router *Router,
	logger Logger,
	metrics MetricsReporter,
	db Database,
	cache Cache,
	mq MessageQueue,
) *Server {
	return &Server{
		router: router,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      router.Handler(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		logger:  logger,
		metrics: metrics,
		db:      db,
		cache:   cache,
		mq:      mq,
		config:  cfg,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting http server", "addr", s.server.Addr)
	
	// 启动服务器指标收集
	s.metrics.IncrementCounter("http_server_start")
	
	// 启动健康检查
	go s.startHealthCheck()
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("failed to start server", "error", err)
		s.metrics.IncrementCounter("http_server_error")
		return err
	}
	
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping http server")
	
	// 记录关闭指标
	s.metrics.IncrementCounter("http_server_stop")
	
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("failed to stop server", "error", err)
		s.metrics.IncrementCounter("http_server_error")
		return err
	}
	
	return nil
}

func (s *Server) startHealthCheck() {
	if !s.config.Health.Enabled {
		return
	}

	// 等待初始延迟
	time.Sleep(s.config.Health.InitialDelay)

	checker := &healthChecker{
		config:  s.config.Health,
		logger:  s.logger,
		metrics: s.metrics,
		checks: []Check{
			{Name: "database", Check: s.db.Ping, Required: true},
			{Name: "cache", Check: s.cache.Ping, Required: false},
			{Name: "mq", Check: s.mq.Ping, Required: false},
		},
	}

	go checker.start()
}

func (s *Server) checkHealth() error {
	// 添加上下文和超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查数据库连接
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database check failed: %w", err)
	}

	// 检查缓存连接
	if err := s.cache.Ping(ctx); err != nil {
		return fmt.Errorf("cache check failed: %w", err)
	}

	// 检查消息队列连接
	if err := s.mq.Ping(ctx); err != nil {
		return fmt.Errorf("message queue check failed: %w", err)
	}

	return nil
}

func (h *healthChecker) start() {
	ticker := time.NewTicker(h.config.Interval)
	defer ticker.Stop()

	done := make(chan struct{})
	defer close(done)

	for {
		select {
		case <-ticker.C:
			h.runChecks()
		case <-done:
			return
		}
	}
}

func (h *healthChecker) runChecks() {
	h.mu.Lock()
	defer h.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), h.config.Timeout)
	defer cancel()

	for _, check := range h.checks {
		if err := check.Check(ctx); err != nil {
			h.logger.Error("health check failed",
				"check", check.Name,
				"error", err,
			)
			h.metrics.IncrementCounter("health_check_failure", "check", check.Name)
			
			if check.Required {
				h.failures++
				if h.failures >= h.config.FailureThreshold {
					h.logger.Error("health check threshold exceeded",
						"failures", h.failures,
						"threshold", h.config.FailureThreshold,
					)
					// TODO: 触发告警或自动恢复机制
				}
			}
			continue
		}
		
		h.metrics.IncrementCounter("health_check_success", "check", check.Name)
		if check.Required {
			h.failures = 0
		}
	}
} 