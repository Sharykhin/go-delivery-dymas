package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
)

const Topic = "latest_position_courier"
const Partition = 0

type CourierPublisher struct {
	publisher sarama.AsyncProducer
	Topic     string
	Partition int32
}

func (courierPublisher *CourierPublisher) PublisherFactory(config *sarama.Config, address string) error {
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return fmt.Errorf("failed to publish Sarama message: %w", err)
	}

	courierPublisher.publisher = producer
	courierPublisher.Topic = Topic
	courierPublisher.Partition = Partition
	return nil
}

func (courierPublisher *CourierPublisher) PublishLatestCourierGeoPositionMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

type CourierPublisherService struct {
	publisher CourierPublisher
}

func (cs *CourierPublisherService) PublishLastCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return err
	}
	cs.publisher.PublishLatestCourierGeoPositionMessage(sarama.ProducerMessage{
		Topic:     cs.publisher.Topic,
		Partition: cs.publisher.Partition,
		Value:     sarama.StringEncoder(message),
	})

	return nil
}

func NewCourierService(
	publisher CourierPublisher,
) domain.CourierPublisherServiceInterface {
	return &CourierPublisherService{
		publisher: publisher,
	}
}

func NewPublisher(config *sarama.Config, address string) (CourierPublisher, error) {
	publisher := CourierPublisher{}
	publisherLink := &publisher
	err := publisherLink.PublisherFactory(config, address)
	return publisher, err
}
