package output

import "context"

type Tracer interface {
    StartSpan(ctx context.Context, name string) (context.Context, Span)
    Close() error
}

type Span interface {
    End()
    SetTag(key string, value interface{})
} 