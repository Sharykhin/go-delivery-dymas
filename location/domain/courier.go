package domain

import (
	"encoding/json"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	"github.com/Shopify/sarama"
	"time"
)

type CourierServiceInterface interface {
	SendData(data *redis.CourierRepositoryData, topic string) error
}
type MessageKafka struct {
	CourierId string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type CourierService struct {
	publisher kafka.CourierPublisherInterface
	repo      redis.CourierRepositoryInterface
}

func (cs CourierService) SendData(data *redis.CourierRepositoryData, topic string) error {
	cs.repo.SaveLatestCourierGeoPosition(data)
	message, err := json.Marshal(MessageKafka{
		CourierId: data.CourierID,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return err
	}
	cs.publisher.PublishMessage(sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return nil
}

func NewCourierService(
	publisher kafka.CourierPublisherInterface,
	repo redis.CourierRepositoryInterface,
) CourierService {
	return CourierService{
		publisher: publisher,
		repo:      repo,
	}
}
