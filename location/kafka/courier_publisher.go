package kafka

import (
	"fmt"
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

func (courierPublisher *CourierPublisher) PublishMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func NewPublisher(config *sarama.Config, address string) (CourierPublisher, error) {
	publisher := CourierPublisher{}
	publisherLink := &publisher
	err := publisherLink.PublisherFactory(config, address)
	return publisher, err
}
