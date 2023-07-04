package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
)

const topic = "latest_position_courier"

type CourierLocationLatestPublisher struct {
	publisher sarama.AsyncProducer
}

func NewCourierLocationPublisher(address string) (*CourierLocationLatestPublisher, error) {
	courierPublisher := CourierLocationLatestPublisher{}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to New async producer: %w", err)
	}

	courierPublisher.publisher = producer
	return &courierPublisher, nil
}

func (courierPublisher *CourierLocationLatestPublisher) publishLatestCourierGeoPositionMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func (courierPublisher *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}
	courierPublisher.publishLatestCourierGeoPositionMessage(sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}
