package query

import "context"

type Bus interface {
    Execute(ctx context.Context, query interface{}) (interface{}, error)
    Register(queryType interface{}, handler Handler)
}

type Handler interface {
    Handle(ctx context.Context, query interface{}) (interface{}, error)
}

type Cacheable interface {
    CacheKey() string
    TTL() time.Duration
} 