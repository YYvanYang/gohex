CREATE TABLE events (
    id VARCHAR(36) PRIMARY KEY,
    aggregate_id VARCHAR(36) NOT NULL,
    type VARCHAR(100) NOT NULL,
    version INT NOT NULL,
    data JSON NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP NULL,
    UNIQUE KEY uk_aggregate_version (aggregate_id, version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_events_aggregate ON events(aggregate_id);
CREATE INDEX idx_events_type ON events(type); 