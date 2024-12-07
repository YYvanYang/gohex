package output

import (
	"context"
	"github.com/your-org/your-project/internal/domain/event"
)

type EventPublisher interface {
	// Publish 发布事件
	Publish(ctx context.Context, events ...event.Event) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler)
}

type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event event.Event) error
} 