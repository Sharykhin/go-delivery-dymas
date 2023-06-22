package domain

import (
	"context"
	"encoding/json"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Shopify/sarama"
	"time"
)

type PublishLastCourierLocationInterface interface {
	SendData(ctx context.Context, data *CourierLocation) error
}
type CourierLocation struct {
	CourierID string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type PublishLastCourierLocation struct {
	publisher kafka.CourierPublisher
	repo      CourierRepositoryInterface
}

func (cs PublishLastCourierLocation) SendData(ctx context.Context, data *CourierLocation) error {
	cs.repo.SaveLatestCourierGeoPosition(ctx, data)
	message, err := json.Marshal(CourierLocation{
		CourierID: data.CourierID,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return err
	}
	cs.publisher.PublishMessage(sarama.ProducerMessage{
		Topic:     cs.publisher.Topic,
		Partition: cs.publisher.Partition,
		Value:     sarama.StringEncoder(message),
	})

	return nil
}

func NewCourierService(
	publisher kafka.CourierPublisher,
	repo CourierRepositoryInterface,
) PublishLastCourierLocationInterface {
	return PublishLastCourierLocation{
		publisher: publisher,
		repo:      repo,
	}
}

type CourierRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, data *CourierLocation) error
}
