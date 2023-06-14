package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
	"time"
)

type CourierPublisher struct {
	publisher sarama.AsyncProducer
}

type CourierService struct {
	publisher CourierPublisher
	repo      domain.CourierRepositoryInterface
}

func (cs CourierService) SendData(data *domain.CourierRepositoryData, ctx context.Context, topic string, partition int32) error {
	cs.repo.SaveLatestCourierGeoPosition(data, ctx)
	message, err := json.Marshal(domain.MessageKafka{
		CourierId: data.CourierID,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return err
	}
	cs.publisher.PublishMessage(sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(message),
	})

	return nil
}

func NewCourierService(
	publisher CourierPublisher,
	repo domain.CourierRepositoryInterface,
) domain.CourierServiceInterface {
	return CourierService{
		publisher: publisher,
		repo:      repo,
	}
}
func (courierPublisher *CourierPublisher) PublisherFactory(config *sarama.Config, address string) error {
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		return fmt.Errorf("failed to publish Sarama message: %w", err)
	}

	courierPublisher.publisher = producer

	return nil
}

func (courierPublisher *CourierPublisher) PublishMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func NewPublisher(config *sarama.Config, address string) (CourierPublisher, error) {
	publisher := CourierPublisher{}
	publisherLink := &publisher
	err := publisherLink.PublisherFactory(config, address)
	return publisher, err
}
