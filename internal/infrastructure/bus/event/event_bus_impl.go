package event

import (
	"context"
	"sync"
	"github.com/your-org/your-project/internal/domain/event"
	"github.com/your-org/your-project/internal/application/port/output"
)

type eventBus struct {
	handlers map[string][]output.EventHandler
	mu       sync.RWMutex
	logger   Logger
	metrics  MetricsReporter
}

func NewEventBus(logger Logger, metrics MetricsReporter) output.EventBus {
	return &eventBus{
		handlers: make(map[string][]output.EventHandler),
		logger:   logger,
		metrics:  metrics,
	}
}

func (b *eventBus) Publish(ctx context.Context, events ...event.Event) error {
	for _, evt := range events {
		timer := b.metrics.StartTimer("event_publish_duration", "type", evt.Type())
		defer timer.Stop()

		b.mu.RLock()
		handlers := b.handlers[evt.Type()]
		b.mu.RUnlock()

		for _, handler := range handlers {
			if err := handler.Handle(ctx, evt); err != nil {
				b.logger.Error("failed to handle event",
					"type", evt.Type(),
					"error", err,
				)
				b.metrics.IncrementCounter("event_handle_failure", "type", evt.Type())
				return err
			}
			b.metrics.IncrementCounter("event_handle_success", "type", evt.Type())
		}
	}
	return nil
}

func (b *eventBus) Subscribe(eventType string, handler output.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	b.logger.Info("subscribed to event", "type", eventType)
} 