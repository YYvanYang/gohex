package handler

import (
	"context"
	"fmt"
	"github.com/your-org/your-project/internal/domain/event"
)

type UserEventHandler struct {
	cache    port.Cache
	emailSvc port.EmailService
	logger   Logger
	metrics  MetricsReporter
}

func (h *UserEventHandler) Handle(ctx context.Context, evt event.Event) error {
	switch e := evt.(type) {
	case *event.UserCreatedEvent:
		return h.handleUserCreated(ctx, e)
	case *event.UserProfileUpdatedEvent:
		return h.handleProfileUpdated(ctx, e)
	case *event.PasswordChangedEvent:
		return h.handlePasswordChanged(ctx, e)
	default:
		return nil
	}
}

func (h *UserEventHandler) handleUserCreated(ctx context.Context, evt *event.UserCreatedEvent) error {
	// 1. 发送欢迎邮件
	if err := h.emailSvc.SendWelcomeEmail(evt.Email, evt.Name); err != nil {
		h.logger.Error("failed to send welcome email", "error", err)
	}

	// 2. 记录指标
	h.metrics.IncrementCounter("user_created")
	return nil
}

func (h *UserEventHandler) handleProfileUpdated(ctx context.Context, evt *event.UserProfileUpdatedEvent) error {
	// 1. 清除用户缓存
	cacheKey := fmt.Sprintf("user:%s", evt.AggregateID())
	if err := h.cache.Delete(ctx, cacheKey); err != nil {
		h.logger.Error("failed to clear user cache", "error", err)
	}

	// 2. 记录指标
	h.metrics.IncrementCounter("user_profile_updated")
	return nil
}

func (h *UserEventHandler) handlePasswordChanged(ctx context.Context, evt *event.PasswordChangedEvent) error {
	// 1. 发送密码变更通知
	if err := h.emailSvc.SendPasswordChangedNotification(evt.Email); err != nil {
		h.logger.Error("failed to send password changed notification", "error", err)
	}

	// 2. 记录指标
	h.metrics.IncrementCounter("user_password_changed")
	return nil
}

// ... 其他事件处理方法 