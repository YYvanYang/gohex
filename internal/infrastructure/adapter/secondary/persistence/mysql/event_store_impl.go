package mysql

import (
    "context"
    "database/sql"
    "encoding/json"
    "github.com/your-org/your-project/internal/domain/event"
)

type mysqlEventStore struct {
    db      *sql.DB
    logger  Logger
    metrics MetricsReporter
}

func NewMySQLEventStore(db *sql.DB, logger Logger, metrics MetricsReporter) *mysqlEventStore {
    return &mysqlEventStore{
        db:      db,
        logger:  logger,
        metrics: metrics,
    }
}

func (s *mysqlEventStore) SaveEvents(ctx context.Context, aggregateID string, events []event.Event, expectedVersion int) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 检查版本
    var currentVersion int
    err = tx.QueryRowContext(ctx, 
        "SELECT MAX(version) FROM events WHERE aggregate_id = ?", 
        aggregateID,
    ).Scan(&currentVersion)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    if currentVersion != expectedVersion {
        return errors.ErrConcurrencyConflict
    }

    // 保存事件
    for _, evt := range events {
        data, err := json.Marshal(evt)
        if err != nil {
            return err
        }

        _, err = tx.ExecContext(ctx, `
            INSERT INTO events (
                id, aggregate_id, type, version, data, occurred_at
            ) VALUES (?, ?, ?, ?, ?, ?)
        `,
            uuid.New().String(),
            evt.AggregateID(),
            evt.Type(),
            evt.Version(),
            data,
            evt.OccurredAt(),
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (s *mysqlEventStore) GetEvents(ctx context.Context, aggregateID string) ([]event.Event, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT type, version, data, occurred_at 
        FROM events 
        WHERE aggregate_id = ? 
        ORDER BY version ASC
    `, aggregateID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []event.Event
    for rows.Next() {
        var (
            eventType   string
            version    int
            data      []byte
            occurredAt time.Time
        )

        if err := rows.Scan(&eventType, &version, &data, &occurredAt); err != nil {
            return nil, err
        }

        evt, err := s.deserializeEvent(eventType, data)
        if err != nil {
            return nil, err
        }
        events = append(events, evt)
    }

    return events, nil
}

func (s *mysqlEventStore) deserializeEvent(eventType string, data []byte) (event.Event, error) {
    switch eventType {
    case event.UserCreated:
        var evt event.UserCreatedEvent
        if err := json.Unmarshal(data, &evt); err != nil {
            return nil, err
        }
        return &evt, nil
    case event.UserProfileUpdated:
        var evt event.UserProfileUpdatedEvent
        if err := json.Unmarshal(data, &evt); err != nil {
            return nil, err
        }
        return &evt, nil
    // ... 其他事件类型
    default:
        return nil, errors.ErrUnknownEventType
    }
} 