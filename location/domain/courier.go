package domain

import (
	"context"
	"fmt"
	"time"
)

type CourierLocationServiceInterface interface {
	SaveLatestCourierLocation(
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

type CourierLocationRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *CourierLocation) error
}

type CourierLocationPublisherInterface interface {
	PublishLatestCourierLocation(ctx context.Context, courierLocation *CourierLocation) error
}

type CourierLocationService struct {
	repo      CourierLocationRepositoryInterface
	publisher CourierLocationPublisherInterface
}

func (courierLocationService *CourierLocationService) SaveLatestCourierLocation(ctx context.Context, courierLocation *CourierLocation) error {
	err := courierLocationService.repo.SaveLatestCourierGeoPosition(ctx, courierLocation)
	if err != nil {
		return fmt.Errorf("failed to store latest courier location in the repository: %w", err)
	}
	err = courierLocationService.publisher.PublishLatestCourierLocation(ctx, courierLocation)

	if err != nil {
		return fmt.Errorf("failed to publish latest courier location: %w", err)
	}

	return nil
}

func NewCourierLocation(id string, latitude, longitude float64) *CourierLocation {
	return &CourierLocation{
		CourierID: id,
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}
}

func NewCourierLocationService(repo CourierLocationRepositoryInterface, courierLocationPublisher CourierLocationPublisherInterface) *CourierLocationService {
	return &CourierLocationService{
		repo:      repo,
		publisher: courierLocationPublisher,
	}
}
