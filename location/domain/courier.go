package domain

import (
	"context"
	"encoding/json"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Shopify/sarama"
	"time"
)

type CourierPublisherServiceInterface interface {
	PublishLastCourierLocation(ctx context.Context, courierLocation *CourierLocation) error
}
type CourierLocation struct {
	CourierID string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type CourierPublisherService struct {
	publisher kafka.CourierPublisher
	repo      CourierRepositoryInterface
}

func (cs *CourierPublisherService) PublishLastCourierLocation(ctx context.Context, courierLocation *CourierLocation) error {
	err := cs.repo.SaveLatestCourierGeoPosition(ctx, courierLocation)
	if err != nil {
		return err
	}
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
	publisher kafka.CourierPublisher,
	repo CourierRepositoryInterface,
) CourierPublisherServiceInterface {
	return &CourierPublisherService{
		publisher: publisher,
		repo:      repo,
	}
}

type CourierRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *CourierLocation) error
}

func CourierLocationFactory(id string, latitude, longitude float64) *CourierLocation {
	return &CourierLocation{
		CourierID: id,
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}
}
