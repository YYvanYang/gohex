package output

import (
	"context"
	"github.com/your-org/your-project/internal/domain/event"
	"time"
)

type EventStore interface {
	SaveEvents(ctx context.Context, aggregateID string, events []event.Event, expectedVersion int) error
	GetEvents(ctx context.Context, aggregateID string) ([]event.Event, error)
	GetEventsFrom(ctx context.Context, aggregateID string, fromVersion int) ([]event.Event, error)
	GetEventsByType(ctx context.Context, eventType string) ([]event.Event, error)
	GetEventsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]event.Event, error)
	GetAggregateHistory(ctx context.Context, aggregateID string) (*AggregateHistory, error)
}

type AggregateHistory struct {
	AggregateID string
	Events      []event.Event
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
} 