package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrCourierNotFound = errors.New("courier was not found")

// CourierLocationServiceInterface saves courier position in storage.
type CourierLocationServiceInterface interface {
	SaveLatestCourierLocation(
		ctx context.Context,
		courierLocation *CourierLocation,
	) error
}

type WorkerLocation interface {
	Init()
	AddTask(courierLocation CourierLocation)
}

// CourierLocation provides information about coords courier.
type CourierLocation struct {
	CourierID string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

// CourierLocationRepositoryInterface saves latest location position courier in storage.
type CourierLocationRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *CourierLocation) error
}

// CourierRepositoryInterface gets latest position by uuid from storage.
type CourierRepositoryInterface interface {
	CourierLocationRepositoryInterface
	GetLatestPositionCourierByID(ctx context.Context, courierID string) (*CourierLocation, error)
}

// CourierLocationPublisherInterface publish message some systems.
type CourierLocationPublisherInterface interface {
	PublishLatestCourierLocation(ctx context.Context, courierLocation *CourierLocation) error
}

// CourierLocationService saves and publishes courier location.
type CourierLocationService struct {
	repo      CourierLocationRepositoryInterface
	publisher CourierLocationPublisherInterface
}

// SaveLatestCourierLocation saves and publish courier location.
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

// NewCourierLocation creates model currier location with current data.
func NewCourierLocation(id string, latitude, longitude float64) CourierLocation {
	return CourierLocation{
		CourierID: id,
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}
}

// NewCourierLocationService creates courier service and init.
func NewCourierLocationService(repo CourierLocationRepositoryInterface, courierLocationPublisher CourierLocationPublisherInterface) *CourierLocationService {
	return &CourierLocationService{
		repo:      repo,
		publisher: courierLocationPublisher,
	}
}
