package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
)

const topic = "latest_position_courier"

type Publisher struct {
	producer sarama.AsyncProducer
}

type HandlerMessageJson interface {
	HandleJsonMessage(ctx context.Context) ([]byte, error)
}

func NewCourierLocationPublisher(address string) (*CourierLocationLatestPublisher, error) {
	publisher := Publisher{}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new sarama async producer: %w", err)
	}

	publisher.producer = producer
	return &publisher, nil
}

func (publisher *Publisher) publish(message sarama.ProducerMessage) {
	publisher.producer.Input() <- &message
}

func (publisher *Publisher) PublishMessage(message []byte) error {
	publisher.publish(sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}
