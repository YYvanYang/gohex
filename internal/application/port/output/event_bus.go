package output

import (
	"context"
	"github.com/your-org/your-project/internal/domain/event"
)

type EventBus interface {
	Publish(ctx context.Context, events ...event.Event) error
	Subscribe(eventType string, handler EventHandler)
	Unsubscribe(eventType string, handler EventHandler)
	Close() error
}

type EventHandler interface {
	Handle(ctx context.Context, event event.Event) error
	HandlerID() string
} 