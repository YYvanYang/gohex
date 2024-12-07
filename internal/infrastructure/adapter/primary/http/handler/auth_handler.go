package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/your-org/your-project/internal/application/command"
	"github.com/your-org/your-project/internal/application/dto"
)

type AuthHandler struct {
	commandBus  command.Bus
	queryBus    query.Bus
	authService AuthService
	logger      Logger
}

func NewAuthHandler(
	commandBus command.Bus,
	queryBus query.Bus,
	authService AuthService,
	logger Logger,
) *AuthHandler {
	return &AuthHandler{
		commandBus:  commandBus,
		queryBus:    queryBus,
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	// 1. 绑定请求
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 2. 执行登录命令
	cmd := &command.LoginCommand{
		Email:     req.Email,
		Password:  req.Password,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	result, err := h.commandBus.Dispatch(c.Request().Context(), cmd)
	if err != nil {
		h.logger.Error("login failed", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// Logout 处理用户登出
func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	token := extractToken(c.Request().Header.Get("Authorization"))
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
	}

	cmd := &command.LogoutCommand{
		UserID: userID,
		Token:  token,
	}

	if err := h.commandBus.Dispatch(c.Request().Context(), cmd); err != nil {
		h.logger.Error("logout failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// 1. 绑定请求
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 2. 执行刷新命令
	cmd := &command.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}
	result, err := h.commandBus.Dispatch(c.Request().Context(), cmd)
	if err != nil {
		h.logger.Error("refresh token failed", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Register(c echo.Context) error {
	// 1. 绑定请求
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Name     string `json:"name" validate:"required"`
		Bio      string `json:"bio" validate:"max=500"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 2. 验证请求
	if err := h.validator.Struct(req); err != nil {
		return h.handleValidationError(err)
	}

	// 3. 执行注册命令
	cmd := &command.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Bio:      req.Bio,
	}

	result, err := h.commandBus.Dispatch(c.Request().Context(), cmd)
	if err != nil {
		h.logger.Error("registration failed", "error", err)
		return h.handleError(err)
	}

	// 4. 记录指标
	h.metrics.IncrementCounter("user_registered")

	return c.JSON(http.StatusCreated, result)
} 