package middleware

import (
	"github.com/gohex/gohex/pkg/tracer"
	"github.com/gohex/gohex/pkg/logger"
	"github.com/gohex/gohex/pkg/metrics"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type TracingMiddleware struct {
	tracer  Tracer
	logger  Logger
	metrics MetricsReporter
}

func NewTracingMiddleware(tracer Tracer, logger Logger, metrics MetricsReporter) *TracingMiddleware {
	return &TracingMiddleware{
		tracer:  tracer,
		logger:  logger,
		metrics: metrics,
	}
}

func (m *TracingMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		
		// 创建或继承 span
		spanCtx, span := m.tracer.StartSpan(req.Context(), "http_request",
			tracer.Tag("http.method", req.Method),
			tracer.Tag("http.url", req.URL.String()),
		)
		defer span.End()

		// 注入跟踪头
		m.tracer.Inject(spanCtx, req.Header)

		// 设置请求上下文
		c.SetRequest(req.WithContext(spanCtx))

		// 执行请求
		err := next(c)

		// 记录响应信息
		if err != nil {
			span.SetTag("error", true)
			span.SetTag("error.message", err.Error())
		}
		span.SetTag("http.status_code", c.Response().Status)

		return err
	}
}

func Tracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := opentracing.StartSpan(
				"gohex.http.request",
				opentracing.Tag{Key: "service", Value: "gohex"},
			)
			defer span.Finish()
			
			// ...
		}
	}
} 