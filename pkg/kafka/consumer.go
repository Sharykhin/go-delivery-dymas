package kafka

import (
	"context"
	"github.com/IBM/sarama"
)

type JSONHandler interface {
	HandleJSONMessage(ctx context.Context, payload any) error
}

type Consumer struct {
	verboseEnabled bool
	consumerGroup  *sarama.ConsumerGroup
}

func WithVerboseConsumer(isEnabled bool) func(*Consumer) {
	return func(consumer *Consumer) {
		consumer.verboseEnabled = isEnabled
	}
}

func NewConsumer(opts ...func(*Consumer)) *Consumer {
	consumer := &Consumer{}
	// create consumer group
	for _, opt := range opts {
		opt(consumer)
	}

	consumerGroup, err := sarama.NewConsumerGroup(...)

	return consumer
}

func (c *Consumer) RegisterJSONHandler(ctx context.Context, topic, handler JSONHandler) {

}

