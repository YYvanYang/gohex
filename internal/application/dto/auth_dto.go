package dto

import (
	"time"
	"github.com/gohex/gohex/internal/domain/vo"
)

// LoginRequestDTO 登录请求
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponseDTO 登录响应
type LoginResponseDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// RefreshTokenRequestDTO 刷新令牌请求
type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponseDTO 刷新令牌响应
type RefreshTokenResponseDTO struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// TokenInfoDTO 令牌信息
type TokenInfoDTO struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AuthErrorDTO 认证错误
type AuthErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RegisterRequestDTO 注册请求
type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"max=500"`
}

// RegisterResponseDTO 注册响应
type RegisterResponseDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ChangePasswordRequestDTO 修改密码请求
type ChangePasswordRequestDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequestDTO 重置密码请求
type ResetPasswordRequestDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// 添加 DTO 到 VO 的转换方法
func (r *RegisterRequestDTO) ToValueObjects() (vo.Email, vo.Password, vo.UserProfile, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return vo.Email{}, vo.Password{}, vo.UserProfile{}, err
	}

	password, err := vo.NewPassword(r.Password)
	if err != nil {
		return vo.Email{}, vo.Password{}, vo.UserProfile{}, err
	}

	profile, err := vo.NewUserProfile(r.Name, r.Bio)
	if err != nil {
		return vo.Email{}, vo.Password{}, vo.UserProfile{}, err
	}

	return email, password, profile, nil
} 