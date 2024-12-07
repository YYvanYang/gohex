package http

import (
	"github.com/labstack/echo/v4"
	"github.com/your-org/your-project/internal/infrastructure/adapter/primary/http/handler"
	"github.com/your-org/your-project/internal/infrastructure/adapter/primary/http/middleware"
	"net/http"
	"fmt"
)

type Router struct {
	echo    *echo.Echo
	logger  Logger
	metrics MetricsReporter
	config  APIConfig
}

type APIConfig struct {
	Version     string `yaml:"version"`
	Deprecated  bool   `yaml:"deprecated"`
	SunsetDate  string `yaml:"sunset_date"`
}

func NewRouter(
	logger Logger,
	metrics MetricsReporter,
	commandBus command.Bus,
	queryBus query.Bus,
	authService AuthService,
) *Router {
	e := echo.New()
	
	// 自定义错误处理
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{} = "Internal Server Error"
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message
		} else if ae, ok := err.(*errors.AppError); ok {
			code = ae.HTTPStatusCode()
			msg = ae.Error()
		}

		// 记录错误
		if code >= 500 {
			logger.Error("request failed", 
				"path", c.Request().URL.Path,
				"error", err,
			)
		}

		// 不在生产环境暴露内部错误
		if code == http.StatusInternalServerError && config.IsProduction() {
			msg = "Internal Server Error"
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				c.NoContent(code)
			} else {
				c.JSON(code, map[string]interface{}{
					"error": msg,
				})
			}
		}
	}
	
	// API 版本
	v1 := e.Group("/api/v1")
	
	// API 文档
	if cfg.Swagger.Enabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	
	// 健康检查
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	
	// 全局中间件
	e.Use(middleware.NewRecoveryMiddleware(logger, metrics))
	e.Use(middleware.NewLoggerMiddleware(logger))
	e.Use(middleware.NewMetricsMiddleware(metrics))
	e.Use(middleware.NewTracingMiddleware())
	
	// 创建处理器
	authHandler := handler.NewAuthHandler(commandBus, queryBus, authService, logger)
	userHandler := handler.NewUserHandler(commandBus, queryBus, logger)
	
	// 认证路由
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout, middleware.RequireAuth(authService))
	}
	
	// 用户路由
	users := v1.Group("/users", middleware.RequireAuth(authService))
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.PUT("/:id/status", userHandler.UpdateUserStatus)
		users.PUT("/:id/password", userHandler.ChangePassword)
	}
	
	return &Router{
		echo:    e,
		logger:  logger,
		metrics: metrics,
	}
}

func (r *Router) Handler() http.Handler {
	return r.echo
}

func (r *Router) addVersionHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 添加版本头
		c.Response().Header().Set("API-Version", r.config.API.Version)
		
		// 如果 API 已弃用，添加相关头
		if r.config.API.Deprecated {
			c.Response().Header().Set("Deprecation", "true")
			if r.config.API.SunsetDate != "" {
				c.Response().Header().Set("Sunset", r.config.API.SunsetDate)
			}
		}
		
		return next(c)
	}
}

// 添加错误处理中间件
func (r *Router) errorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  = "Internal Server Error"
		details map[string]interface{} = nil
	)

	switch e := err.(type) {
	case *echo.HTTPError:
		code = e.Code
		msg = fmt.Sprintf("%v", e.Message)
	case *errors.AppError:
		code = e.HTTPStatusCode()
		msg = e.Message
		details = e.Details
	case validator.ValidationErrors:
		code = http.StatusBadRequest
		msg = "Validation Error"
		details = make(map[string]interface{})
		for _, err := range e {
			details[err.Field()] = err.Tag()
		}
	}

	// 记录错误
	if code >= 500 {
		r.logger.Error("request failed",
			"path", c.Request().URL.Path,
			"method", c.Request().Method,
			"error", err,
			"code", code,
		)
		r.metrics.IncrementCounter("http_error_5xx")
	} else {
		r.logger.Warn("request failed",
			"path", c.Request().URL.Path,
			"method", c.Request().Method,
			"error", err,
			"code", code,
		)
		r.metrics.IncrementCounter("http_error_4xx")
	}

	// 构造响应
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"message": msg,
			"code":    code,
		},
	}

	if details != nil {
		response["error"].(map[string]interface{})["details"] = details
	}

	// 发送响应
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(code)
		} else {
			c.JSON(code, response)
		}
	}
} 