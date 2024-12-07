package middleware

import (
	"strings"
	"github.com/labstack/echo/v4"
	"github.com/your-org/your-project/internal/application/query"
	"net/http"
)

type AuthMiddleware struct {
	queryBus query.QueryBus
	logger   Logger
	metrics  MetricsReporter
}

func NewAuthMiddleware(
	queryBus query.QueryBus,
	logger Logger,
	metrics MetricsReporter,
) *AuthMiddleware {
	return &AuthMiddleware{
		queryBus: queryBus,
		logger:   logger,
		metrics:  metrics,
	}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		timer := m.metrics.StartTimer("auth_middleware_duration")
		defer timer.Stop()

		// 1. 提取令牌
		token := extractToken(c.Request().Header.Get("Authorization"))
		if token == "" {
			m.metrics.IncrementCounter("auth_middleware_missing_token")
			return echo.NewHTTPError(401, "missing token")
		}

		// 2. 验证令牌
		q := &query.ValidateTokenQuery{Token: token}
		result, err := m.queryBus.Execute(c.Request().Context(), q)
		if err != nil {
			m.metrics.IncrementCounter("auth_middleware_invalid_token")
			return echo.NewHTTPError(401, "invalid token")
		}

		// 3. 设置上下文
		authInfo := result.(*dto.AuthInfoDTO)
		c.Set("user_id", authInfo.UserID)
		c.Set("roles", authInfo.Roles)

		m.metrics.IncrementCounter("auth_middleware_success")
		return next(c)
	}
}

func (m *AuthMiddleware) RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取用户角色
			userRoles, ok := c.Get("user_roles").([]string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user roles")
			}

			// 检查是否有所需角色
			hasRole := false
			for _, required := range roles {
				for _, role := range userRoles {
					if role == required {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				m.logger.Warn("access denied",
					"path", c.Request().URL.Path,
					"required_roles", roles,
					"user_roles", userRoles,
				)
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}

func extractToken(auth string) string {
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func hasAnyRole(userRoles []string, requiredRoles []string) bool {
	for _, required := range requiredRoles {
		for _, role := range userRoles {
			if role == required {
				return true
			}
		}
	}
	return false
}

func RequireAuth(authService AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取令牌
			token := extractToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}
			
			// 验证令牌
			claims, err := authService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			
			// 检查令牌是否已吊销
			if authService.IsTokenRevoked(c.Request().Context(), token) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token revoked")
			}
			
			// 设置用户信息到上下文
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_roles", claims.Roles)
			
			return next(c)
		}
	}
}

func extractToken(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	return parts[1]
} 