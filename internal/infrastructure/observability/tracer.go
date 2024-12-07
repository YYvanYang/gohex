package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Tracer struct {
	provider trace.TracerProvider
	tracer   trace.Tracer
}

func NewTracer(serviceName, endpoint string) (*Tracer, error) {
	// 创建 Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}

	// 创建 TracerProvider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)

	return &Tracer{
		provider: provider,
		tracer:   provider.Tracer(serviceName),
	}, nil
}

func (t *Tracer) StartSpan(ctx context.Context, name string) (trace.Span, context.Context) {
	return t.tracer.Start(ctx, name)
}

func (t *Tracer) Close() error {
	if provider, ok := t.provider.(*sdktrace.TracerProvider); ok {
		return provider.Shutdown(context.Background())
	}
	return nil
} 