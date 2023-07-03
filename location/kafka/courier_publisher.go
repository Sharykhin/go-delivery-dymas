package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Shopify/sarama"
)

const topic = "latest_position_courier"
const partition = 0

type CourierLocationLatestPublisher struct {
	publisher sarama.AsyncProducer
	topic     string
	partition int32
}

func PublisherCourierLocationFactory(address string) (domain.CourierLocationPublisherInterface, error) {
	courierPublisher := CourierLocationLatestPublisher{}
	config := sarama.NewConfig()
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

func (courierPublisher CourierLocationLatestPublisher) publishLatestCourierGeoPositionMessage(message sarama.ProducerMessage) {
	courierPublisher.publisher.Input() <- &message
}

func (courierPublisher CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return err
	}
	courierPublisher.publishLatestCourierGeoPositionMessage(sarama.ProducerMessage{
		Topic:     courierPublisher.topic,
		Partition: courierPublisher.partition,
		Value:     sarama.StringEncoder(message),
	})

	return nil
}
