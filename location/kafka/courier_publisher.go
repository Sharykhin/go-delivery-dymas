package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
	"time"
)

const topic = "latest_position_courier"
const partition = 0

type CourierPublisher struct {
	publisher sarama.AsyncProducer
}

type CourierService struct {
	publisher CourierPublisher
}

func (cs CourierService) SendData(data *domain.CourierLocationEvent) error {
	message, err := json.Marshal(domain.CourierLocationEvent{
		CourierID: data.CourierID,
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
) domain.CourierServiceInterface {
	return CourierService{
		publisher: publisher,
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
