package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type CourierPublisherInterface interface {
	PublisherFactory(config *sarama.Config, address string)
	PublishMessage(message sarama.ProducerMessage)
}

type CourierPublisher struct {
	publisher sarama.AsyncProducer
}

func (courierPublisher *CourierPublisher) PublisherFactory(config *sarama.Config, address string) {
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, error := sarama.NewAsyncProducer([]string{address}, config)
	if error != nil {
		panic(error)
	}

	courierPublisher.publisher = producer
}

func (courierPublisher *CourierPublisher) PublishMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func NewPublisher(config *sarama.Config, address string) CourierPublisherInterface {
	publisher := &CourierPublisher{}
	publisher.PublisherFactory(config, address)
	fmt.Println(publisher)
	return publisher
}
