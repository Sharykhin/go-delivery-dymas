package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

type CourierLocationLatestPublisher struct {
	publisher sarama.AsyncProducer
	topic     string
}

func NewCourierLocationPublisher(address string, topic string) (*CourierLocationLatestPublisher, error) {
	courierPublisher := CourierLocationLatestPublisher{}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new sarama async producer: %w", err)
	}

	courierPublisher.publisher = producer
	courierPublisher.topic = topic
	return &courierPublisher, nil
}

func (courierPublisher *CourierLocationLatestPublisher) publishLatestCourierGeoPositionMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func (courierPublisher *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *any) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}
	courierPublisher.publishLatestCourierGeoPositionMessage(sarama.ProducerMessage{
		Topic: courierPublisher.topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}
