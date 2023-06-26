package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
)

const topic = "latest_position_courier"
const partition = 0

type CourierPublisher struct {
	publisher sarama.AsyncProducer
	topic     string
	partition int32
}

func PublisherCourierLocationFactory(config *sarama.Config, address string) (CourierPublisher, error) {
	courierPublisher := CourierPublisher{}
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		err = fmt.Errorf("failed to publish Sarama message: %w", err)
	}

	courierPublisher.publisher = producer
	courierPublisher.topic = topic
	courierPublisher.partition = partition
	return courierPublisher, err
}

func (courierPublisher *CourierPublisher) PublishLatestCourierGeoPositionMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

type CourierPublisherService struct {
	publisher CourierPublisher
}

func (cs *CourierPublisherService) PublishLatestCourierLocation(courierLocation *domain.CourierLocation) error {
	fmt.Println(" cs.publisher.topic")
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return err
	}
	cs.publisher.PublishLatestCourierGeoPositionMessage(sarama.ProducerMessage{
		Topic:     cs.publisher.topic,
		Partition: cs.publisher.partition,
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
