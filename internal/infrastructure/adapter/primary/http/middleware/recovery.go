package middleware

import (
	"fmt"
	"runtime"
	"github.com/labstack/echo/v4"
)

type RecoveryMiddleware struct {
	logger  Logger
	metrics MetricsReporter
}

func NewRecoveryMiddleware(logger Logger, metrics MetricsReporter) echo.MiddlewareFunc {
	r := &RecoveryMiddleware{
		logger:  logger,
		metrics: metrics,
	}
	return r.Handle
}

func (m *RecoveryMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack := make([]byte, 4<<10) // 4 KB
				length := runtime.Stack(stack, false)

				m.logger.Error("panic recovered",
					"error", err,
					"stack", string(stack[:length]),
					"url", c.Request().URL.String(),
					"method", c.Request().Method,
				)

				m.metrics.IncrementCounter("panic_recovered",
					"path", c.Request().URL.Path,
				)

				c.Error(echo.NewHTTPError(500, "Internal Server Error"))
			}
		}()
		return next(c)
	}
} 