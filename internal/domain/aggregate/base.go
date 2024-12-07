package aggregate

import (
	"github.com/your-org/your-project/internal/domain/event"
)

// BaseAggregate 提供聚合根的基础实现
type BaseAggregate struct {
	id      string
	version int
	events  []event.Event
}

func NewBaseAggregate(id string) *BaseAggregate {
	return &BaseAggregate{
		id:      id,
		version: 0,
		events:  make([]event.Event, 0),
	}
}

func (a *BaseAggregate) ID() string {
	return a.id
}

func (a *BaseAggregate) Version() int {
	return a.version
}

func (a *BaseAggregate) Events() []event.Event {
	return a.events
}

func (a *BaseAggregate) ClearEvents() {
	a.events = make([]event.Event, 0)
}

func (a *BaseAggregate) AddEvent(event event.Event) {
	a.events = append(a.events, event)
	a.version++
} 