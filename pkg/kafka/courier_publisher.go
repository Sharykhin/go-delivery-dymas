package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

type Publisher struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewPublisher(address string, topic string, Publisher *any) (*any, error) {
	publisher := Publisher
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new sarama async producer: %w", err)
	}

	publisher.producer = producer
	publisher.topic = topic
	return publisher, nil
}

func (publisher *Publisher) publishMessage(message sarama.ProducerMessage) {
	publisher.producer.Input() <- &message
}

func (publisher *Publisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *any) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}
	publisher.publishMessage(sarama.ProducerMessage{
		Topic: publisher.topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}
