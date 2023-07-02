package domain

import (
	"context"
	"time"
)

type CourierPublisherServiceInterface interface {
	PublishLatestCourierLocation(
		ctx context.Context,
		courierLocation *CourierLocation,
	) error
}
type CourierLocation struct {
	CourierID string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type CourierRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *CourierLocation) error
}

type CourierLocationPublisherInterface interface {
	PublishLatestCourierLocation(courierLocation *CourierLocation) error
}

type CourierPublisherService struct {
	repo      CourierRepositoryInterface
	publisher CourierLocationPublisherInterface
}

func (courierPublisherService CourierPublisherService) PublishLatestCourierLocation(ctx context.Context, courierLocation *CourierLocation) error {
	err := courierPublisherService.repo.SaveLatestCourierGeoPosition(ctx, courierLocation)
	if err != nil {
		return err
	}
	err = courierPublisherService.publisher.PublishLatestCourierLocation(courierLocation)

	return err
}

func CourierLocationFactory(id string, latitude, longitude float64) *CourierLocation {
	return &CourierLocation{
		CourierID: id,
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}
}

func CourierPublisherServiceFactory(repo CourierRepositoryInterface, courierLocationPublisher CourierLocationPublisherInterface) CourierPublisherServiceInterface {
	return CourierPublisherService{
		repo:      repo,
		publisher: courierLocationPublisher,
	}
}
