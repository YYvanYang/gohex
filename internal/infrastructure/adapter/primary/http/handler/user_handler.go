package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gohex/gohex/internal/application/dto"
	"github.com/gohex/gohex/internal/application/command"
	"github.com/gohex/gohex/internal/application/query"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/pkg/tracer"
)

// UserHandler 处理所有用户相关的 HTTP 请求
type UserHandler struct {
	commandBus command.Bus
	queryBus   query.Bus
	logger     Logger
	metrics    MetricsReporter
	validator  *validator.Validate
}

// NewUserHandler 创建新的 UserHandler 实例
func NewUserHandler(
	commandBus command.Bus,
	queryBus query.Bus,
	logger Logger,
	metrics MetricsReporter,
) *UserHandler {
	return &UserHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
		metrics:    metrics,
		validator:  validator.New(),
	}
}

// RegisterUserRequest 用户注册请求
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"max=500"`
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c echo.Context) error {
	span, ctx := tracer.StartSpan(c.Request().Context(), "UserHandler.Register")
	defer span.End()

	timer := h.metrics.StartTimer("http_register_user")
	defer timer.Stop()

	// 1. 绑定请求
	var req RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return h.handleError(err)
	}

	// 2. 验证请求
	if err := h.validator.Struct(req); err != nil {
		return h.handleValidationError(err)
	}

	// 3. 构造命令
	cmd := command.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Bio:      req.Bio,
	}

	// 4. 执行命令
	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return h.handleError(err)
	}

	// 5. 返回响应
	return c.JSON(http.StatusCreated, result)
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c echo.Context) error {
	span, ctx := tracer.StartSpan(c.Request().Context(), "UserHandler.GetUser")
	defer span.End()

	userID := c.Param("id")
	if userID == "" {
		return h.handleError(errors.NewValidationError("user_id is required"))
	}

	q := query.GetUserQuery{ID: userID}
	user, err := h.queryBus.Execute(ctx, q)
	if err != nil {
		return h.handleError(err)
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Name string `json:"name" validate:"required"`
	Bio  string `json:"bio" validate:"max=500"`
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	span, ctx := tracer.StartSpan(c.Request().Context(), "UserHandler.UpdateProfile")
	defer span.End()

	userID := c.Param("id")
	if userID == "" {
		return h.handleError(errors.NewValidationError("user_id is required"))
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return h.handleError(err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.handleValidationError(err)
	}

	cmd := command.UpdateProfileCommand{
		UserID: userID,
		Name:   req.Name,
		Bio:    req.Bio,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		return h.handleError(err)
	}

	return c.NoContent(http.StatusOK)
}

// handleError 统一错误处理
func (h *UserHandler) handleError(err error) error {
	h.logger.Error("request failed", "error", err)
	h.metrics.IncrementCounter("http_error")

	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		return echo.NewHTTPError(h.statusCodeFromError(appErr), appErr)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
}

// handleValidationError 处理验证错误
func (h *UserHandler) handleValidationError(err error) error {
	h.metrics.IncrementCounter("validation_error")
	return echo.NewHTTPError(http.StatusBadRequest, errors.NewValidationError(err.Error()))
}

// statusCodeFromError 根据错误类型返回对应的 HTTP 状态码
func (h *UserHandler) statusCodeFromError(err *errors.AppError) int {
	switch err.Code {
	case errors.ErrCodeValidation:
		return http.StatusBadRequest
	case errors.ErrCodeNotFound:
		return http.StatusNotFound
	case errors.ErrCodeConflict:
		return http.StatusConflict
	case errors.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var params struct {
		Page     int    `query:"page" validate:"min=1"`
		PageSize int    `query:"page_size" validate:"min=1,max=100"`
		Status   string `query:"status" validate:"omitempty,oneof=active inactive suspended deleted"`
		SortBy   string `query:"sort_by" validate:"omitempty,oneof=created_at updated_at name email"`
		SortDir  string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
	}

	if err := c.Bind(&params); err != nil {
		return h.handleValidationError(err)
	}

	query := &query.ListUsersQuery{
		Page:     params.Page,
		PageSize: params.PageSize,
		Status:   params.Status,
		SortBy:   params.SortBy,
		SortDir:  params.SortDir,
	}

	result, err := h.queryBus.Execute(ctx, query)
	if err != nil {
		return h.handleError(err)
	}

	return c.JSON(http.StatusOK, result)
} 