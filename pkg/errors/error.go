package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 定义错误码类型
type ErrorCode string

const (
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

// AppError 定义应用错误结构
type AppError struct {
	Code     ErrorCode         `json:"code"`
	Message  string           `json:"message"`
	Details  map[string]interface{}   `json:"details,omitempty"`
	Internal error            `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// HTTPStatusCode 获取对应的 HTTP 状态码
func (e *AppError) HTTPStatusCode() int {
	switch e.Code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// WithInternal 添加内部错误
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func GetErrorCode(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrCodeInternal
} 