package event

import (
	"time"
)

type Event interface {
	AggregateID() string
	Type() string
	OccurredAt() time.Time
	Data() interface{}
}

type BaseEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	data        interface{}
}

func NewBaseEvent(aggregateID, eventType string, data interface{}) BaseEvent {
	return BaseEvent{
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
		data:        data,
	}
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseEvent) Type() string {
	return e.eventType
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseEvent) Data() interface{} {
	return e.data
} 