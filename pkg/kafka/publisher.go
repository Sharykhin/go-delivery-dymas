package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

type Publisher struct {
	producer sarama.AsyncProducer
	topic    string
}

type JSONMessagePublishHandler interface {
	JSONMessageHandle(ctx context.Context) ([]byte, error)
}

func NewPublisher(address, topic string) (*Publisher, error) {
	publisher := Publisher{}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)

	if err != nil {
		return nil, fmt.Errorf("failed to create a new sarama async producer: %w", err)
	}

	publisher.producer = producer
	publisher.topic = topic

	return &publisher, nil
}

func (publisher *Publisher) publish(message sarama.ProducerMessage) {
	publisher.producer.Input() <- &message
}

func (publisher *Publisher) PublishMessage(ctx context.Context, message []byte) error {
	publisher.publish(sarama.ProducerMessage{
		Topic: publisher.topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}
