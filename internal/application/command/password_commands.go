package command

import (
	"context"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/errors"
)

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type ChangePasswordHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	changeCmd := cmd.(*ChangePasswordCommand)

	return nil, h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 获取用户
		user, err := h.userRepo.FindByID(ctx, changeCmd.UserID)
		if err != nil {
			return err
		}

		// 2. 创建密码值对象
		currentPassword, err := vo.NewPassword(changeCmd.CurrentPassword)
		if err != nil {
			return err
		}

		newPassword, err := vo.NewPassword(changeCmd.NewPassword)
		if err != nil {
			return err
		}

		// 3. 修改密码
		if err := user.ChangePassword(currentPassword, newPassword); err != nil {
			return err
		}

		// 4. 保存用户
		if err := h.userRepo.Update(ctx, user); err != nil {
			return err
		}

		// 5. 保存事件
		return h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), user.Version())
	})
}

// ResetPasswordCommand 重置密码命令
type ResetPasswordCommand struct {
	UserID      string
	NewPassword string
}

type ResetPasswordHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *ResetPasswordHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	resetCmd := cmd.(*ResetPasswordCommand)

	return nil, h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 获取用户
		user, err := h.userRepo.FindByID(ctx, resetCmd.UserID)
		if err != nil {
			return err
		}

		// 2. 创建新密码值对象
		newPassword, err := vo.NewPassword(resetCmd.NewPassword)
		if err != nil {
			return err
		}

		// 3. 重置密码
		if err := user.ResetPassword(newPassword); err != nil {
			return err
		}

		// 4. 保存用户
		if err := h.userRepo.Update(ctx, user); err != nil {
			return err
		}

		// 5. 保存事件
		return h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), user.Version())
	})
}

// RequestPasswordResetCommand 请求密码重置命令
type RequestPasswordResetCommand struct {
	Email string `validate:"required,email"`
}

type RequestPasswordResetHandler struct {
	userRepo   port.UserRepository
	tokenSvc   port.TokenService
	emailSvc   port.EmailService
	logger     Logger
	metrics    MetricsReporter
}

func (h *RequestPasswordResetHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	resetCmd := cmd.(*RequestPasswordResetCommand)

	// 1. 查找用户
	email, err := vo.NewEmail(resetCmd.Email)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			// 为了安全，即使用户不存在也返回成功
			return nil, nil
		}
		return nil, err
	}

	// 2. 生成重置令牌
	token, err := h.tokenSvc.GeneratePasswordResetToken(user)
	if err != nil {
		return nil, err
	}

	// 3. 发送重置邮件
	if err := h.emailSvc.SendPasswordResetEmail(user.Email().String(), token); err != nil {
		h.logger.Error("failed to send password reset email", "error", err)
		return nil, err
	}

	return nil, nil
} 