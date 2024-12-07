package listener

import (
	"context"
	"github.com/your-org/your-project/internal/domain/event"
)

type UserEventListener struct {
	logger  Logger
	metrics MetricsReporter
	cache   Cache
}

func NewUserEventListener(logger Logger, metrics MetricsReporter, cache Cache) *UserEventListener {
	return &UserEventListener{
		logger:  logger,
		metrics: metrics,
		cache:   cache,
	}
}

func (l *UserEventListener) Handle(ctx context.Context, event event.Event) error {
	span, ctx := tracer.StartSpan(ctx, "UserEventListener.Handle")
	defer span.End()

	timer := l.metrics.StartTimer("event_handler_duration")
	defer timer.Stop()

	switch e := event.(type) {
	case *event.UserCreatedEvent:
		return l.handleUserCreated(ctx, e)
	case *event.UserProfileUpdatedEvent:
		return l.handleProfileUpdated(ctx, e)
	case *event.UserStatusChangedEvent:
		return l.handleStatusChanged(ctx, e)
	case *event.UserLoggedInEvent:
		return l.handleUserLoggedIn(ctx, e)
	default:
		l.logger.Debug("ignoring unknown event type", "type", event.Type())
		return nil
	}
}

func (l *UserEventListener) handleUserCreated(ctx context.Context, e *event.UserCreatedEvent) error {
	l.metrics.IncrementCounter("user_created")
	// 可以在这里发送欢迎邮件等
	return nil
}

func (l *UserEventListener) handleProfileUpdated(ctx context.Context, e *event.UserProfileUpdatedEvent) error {
	// 清除用户缓存
	cacheKey := "user:" + e.AggregateID()
	if err := l.cache.Delete(ctx, cacheKey); err != nil {
		l.logger.Error("failed to clear user cache", "error", err)
	}
	return nil
}

func (l *UserEventListener) handleStatusChanged(ctx context.Context, e *event.UserStatusChangedEvent) error {
	// 处理用户状态变更，比如发送通知等
	return nil
}

func (l *UserEventListener) handleUserLoggedIn(ctx context.Context, e *event.UserLoggedInEvent) error {
	// 记录登录历史，更新最后登录时间等
	return nil
} 