package kafka

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/your-org/your-project/internal/domain/event"
	"github.com/your-org/your-project/internal/application/port"
)

type kafkaEventBus struct {
	producer   sarama.SyncProducer
	consumer   sarama.ConsumerGroup
	handlers   map[string][]port.EventHandler
	logger     Logger
	metrics    MetricsReporter
}

func NewKafkaEventBus(
	producer sarama.SyncProducer,
	consumer sarama.ConsumerGroup,
	logger Logger,
	metrics MetricsReporter,
) port.EventBus {
	return &kafkaEventBus{
		producer:  producer,
		consumer:  consumer,
		handlers:  make(map[string][]port.EventHandler),
		logger:    logger,
		metrics:   metrics,
	}
}

func (b *kafkaEventBus) Publish(ctx context.Context, events ...event.Event) error {
	span, ctx := tracer.StartSpan(ctx, "kafkaEventBus.Publish")
	defer span.End()

	for _, evt := range events {
		timer := b.metrics.StartTimer("event_publish_duration")
		msg, err := b.serializeEvent(evt)
		if err != nil {
			timer.Stop()
			return err
		}

		_, _, err = b.producer.SendMessage(&sarama.ProducerMessage{
			Topic: b.eventToTopic(evt),
			Value: sarama.StringEncoder(msg),
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte("event_type"),
					Value: []byte(evt.Type()),
				},
			},
		})

		timer.Stop()

		if err != nil {
			b.logger.Error("failed to publish event",
				"event_type", evt.Type(),
				"error", err,
			)
			b.metrics.IncrementCounter("event_publish_failure")
			return err
		}

		b.metrics.IncrementCounter("event_publish_success")
	}

	return nil
}

func (b *kafkaEventBus) Subscribe(handler port.EventHandler, eventTypes ...string) {
	for _, eventType := range eventTypes {
		b.handlers[eventType] = append(b.handlers[eventType], handler)
	}
}

func (b *kafkaEventBus) Start(ctx context.Context) error {
	topics := b.getSubscribedTopics()
	
	go func() {
		for {
			err := b.consumer.Consume(ctx, topics, &consumerGroupHandler{
				handlers: b.handlers,
				logger:   b.logger,
				metrics:  b.metrics,
			})
			if err != nil {
				b.logger.Error("failed to consume messages", "error", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

type consumerGroupHandler struct {
	handlers map[string][]port.EventHandler
	logger   Logger
	metrics  MetricsReporter
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		timer := h.metrics.StartTimer("event_process_duration")
		
		eventType := string(msg.Headers[0].Value)
		handlers := h.handlers[eventType]

		event, err := h.deserializeEvent(eventType, msg.Value)
		if err != nil {
			h.logger.Error("failed to deserialize event", "error", err)
			timer.Stop()
			continue
		}

		for _, handler := range handlers {
			if err := handler.Handle(session.Context(), event); err != nil {
				h.logger.Error("failed to handle event",
					"event_type", eventType,
					"error", err,
				)
				h.metrics.IncrementCounter("event_handle_failure")
			} else {
				h.metrics.IncrementCounter("event_handle_success")
			}
		}

		session.MarkMessage(msg, "")
		timer.Stop()
	}
	return nil
}

func (b *kafkaEventBus) serializeEvent(evt event.Event) (string, error) {
	data := map[string]interface{}{
		"id":           evt.AggregateID(),
		"type":         evt.Type(),
		"occurred_at":  evt.OccurredAt(),
		"data":         evt.Data(),
	}
	
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	return string(bytes), nil
}

func (b *kafkaEventBus) eventToTopic(evt event.Event) string {
	return "events." + evt.Type()
}

func (b *kafkaEventBus) getSubscribedTopics() []string {
	topics := make([]string, 0, len(b.handlers))
	for eventType := range b.handlers {
		topics = append(topics, "events."+eventType)
	}
	return topics
} 