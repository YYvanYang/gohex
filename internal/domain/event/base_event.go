package event

import "time"

type Event interface {
    AggregateID() string
    Type() string
    OccurredAt() time.Time
    Version() int
}

type BaseEvent struct {
    aggregateID string
    eventType   string
    occurredAt  time.Time
    version     int
}

func NewBaseEvent(aggregateID string, eventType string) BaseEvent {
    return BaseEvent{
        aggregateID: aggregateID,
        eventType:   eventType,
        occurredAt:  time.Now(),
    }
}

func (e BaseEvent) AggregateID() string { return e.aggregateID }
func (e BaseEvent) Type() string        { return e.eventType }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) Version() int        { return e.version } 