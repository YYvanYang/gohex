package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/your-org/your-project/internal/domain/event"
)

type eventStore struct {
	db      *sql.DB
	logger  Logger
	metrics MetricsReporter
}

type eventModel struct {
	ID           string          `db:"id"`
	AggregateID  string          `db:"aggregate_id"`
	Type         string          `db:"type"`
	Version      int             `db:"version"`
	Data         json.RawMessage `db:"data"`
	OccurredAt   time.Time       `db:"occurred_at"`
	PublishedAt  *time.Time      `db:"published_at"`
}

func NewEventStore(db *sql.DB, logger Logger, metrics MetricsReporter) *eventStore {
	return &eventStore{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}
}

func (s *eventStore) SaveEvents(ctx context.Context, aggregateID string, events []event.Event, expectedVersion int) error {
	span, ctx := tracer.StartSpan(ctx, "eventStore.SaveEvents")
	defer span.End()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 检查版本
	var currentVersion int
	err = tx.QueryRowContext(ctx, 
		"SELECT COALESCE(MAX(version), 0) FROM events WHERE aggregate_id = ?",
		aggregateID,
	).Scan(&currentVersion)

	if err != nil {
		return err
	}

	if currentVersion != expectedVersion {
		return ErrConcurrencyConflict
	}

	// 保存事件
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events (id, aggregate_id, type, version, data, occurred_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, e := range events {
		data, err := json.Marshal(e.Data())
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx,
			uuid.New().String(),
			e.AggregateID(),
			e.Type(),
			expectedVersion+i+1,
			data,
			e.OccurredAt(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *eventStore) GetEvents(ctx context.Context, aggregateID string) ([]event.Event, error) {
	span, ctx := tracer.StartSpan(ctx, "eventStore.GetEvents")
	defer span.End()

	query := `
		SELECT id, aggregate_id, type, version, data, occurred_at
		FROM events
		WHERE aggregate_id = ?
		ORDER BY version ASC
	`

	rows, err := s.db.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var model eventModel
		err := rows.Scan(
			&model.ID,
			&model.AggregateID,
			&model.Type,
			&model.Version,
			&model.Data,
			&model.OccurredAt,
		)
		if err != nil {
			return nil, err
		}

		event, err := s.deserializeEvent(&model)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *eventStore) deserializeEvent(model *eventModel) (event.Event, error) {
	// 根据事件类型反序列化数据
	switch model.Type {
	case event.UserCreated:
		var data event.UserCreatedEvent
		if err := json.Unmarshal(model.Data, &data); err != nil {
			return nil, err
		}
		return &data, nil
	// ... 其他事件类型的处理
	default:
		return nil, ErrUnknownEventType
	}
} 