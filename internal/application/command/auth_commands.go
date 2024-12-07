package command

import (
	"context"
	"time"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/internal/application/port"
)

// LoginCommand 登录命令
type LoginCommand struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	IP        string
	UserAgent string
}

type LoginHandler struct {
	userRepo   port.UserRepository
	tokenSvc   port.TokenService
	eventStore port.EventStore
	cache      port.Cache
	logger     Logger
	metrics    MetricsReporter
}

func (h *LoginHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	loginCmd := cmd.(*LoginCommand)

	// 1. 获取用户
	email, err := vo.NewEmail(loginCmd.Email)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. 验证密码
	if err := user.ValidatePassword(loginCmd.Password); err != nil {
		// 记录失败次数
		h.recordLoginFailure(ctx, user.ID())
		return nil, errors.ErrInvalidCredentials
	}

	// 3. 检查账户状态
	if !user.Status().IsActive() {
		return nil, errors.ErrAccountLocked
	}

	// 4. 生成令牌
	token, expiresAt, err := h.tokenSvc.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// 5. 记录登录事件
	user.RecordLogin(loginCmd.IP, loginCmd.UserAgent)
	if err := h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), user.Version()); err != nil {
		return nil, err
	}

	// 6. 清除失败计数
	h.clearLoginFailures(ctx, user.ID())

	return &dto.LoginResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		RefreshToken: "", // 如果需要刷新令牌，在这里生成
	}, nil
}

func (h *LoginHandler) recordLoginFailure(ctx context.Context, userID string) {
	key := "login_failures:" + userID
	count, _ := h.cache.Increment(ctx, key, 1)
	if count == 1 {
		h.cache.Expire(ctx, key, time.Hour)
	}
	if count >= 5 {
		h.lockAccount(ctx, userID)
	}
}

func (h *LoginHandler) clearLoginFailures(ctx context.Context, userID string) {
	h.cache.Delete(ctx, "login_failures:"+userID)
}

func (h *LoginHandler) lockAccount(ctx context.Context, userID string) {
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		h.logger.Error("failed to lock account", "error", err)
		return
	}

	if err := user.Lock(); err != nil {
		h.logger.Error("failed to lock account", "error", err)
		return
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("failed to save locked account", "error", err)
	}
}

// LogoutCommand 登出命令
type LogoutCommand struct {
	UserID string
	Token  string
}

type LogoutHandler struct {
	tokenSvc port.TokenService
	logger   Logger
	metrics  MetricsReporter
}

func (h *LogoutHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	logoutCmd := cmd.(*LogoutCommand)

	// 1. 吊销令牌
	if err := h.tokenSvc.RevokeToken(ctx, logoutCmd.Token); err != nil {
		h.logger.Error("failed to revoke token", "error", err)
		return nil, err
	}

	return nil, nil
} 