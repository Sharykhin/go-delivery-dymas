package domain

import (
	"context"
	"errors"
	"fmt"
)

var ErrCourierNotFound = errors.New("courier was not found")

// CourierWithLatestPosition Model courier provides information and current latest position
type CourierWithLatestPosition struct {
	LatestPosition *LocationPosition
	FirstName      string
	ID             string
	IsAvailable    bool
}

// CourierClientInterface gets latest location position courier from storage
type CourierClientInterface interface {
	GetLatestPosition(ctx context.Context, courierID string) (*LocationPosition, error)
}

// LocationPosition provides location position courier
type LocationPosition struct {
	Latitude  float64
	Longitude float64
}

// CourierService provides information about courier and latest position from storage
type CourierService struct {
	courierClient     CourierClientInterface
	courierRepository CourierRepositoryInterface
}

// CourierRepositoryInterface saves and reads courier from storage
type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
}

// CourierServiceInterface Gets information about courier and latest position courier from storage
type CourierServiceInterface interface {
	GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierWithLatestPosition, error)
}

// Courier provides information about courier
type Courier struct {
	ID          string
	FirstName   string
	IsAvailable bool
}

// GetCourierWithLatestPosition gets latest position from server and storage
func (s *CourierService) GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierWithLatestPosition, error) {
	var locationPosition *LocationPosition

	courier, err := s.courierRepository.GetCourierByID(
		ctx,
		courierID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get courier from the repository: %w", err)
	}

	courierLatestPositionResponse, err := s.courierClient.GetLatestPosition(ctx, courierID)
	isErrCourierNotFound := errors.Is(err, ErrCourierNotFound)

	if err != nil && !isErrCourierNotFound {
		return nil, fmt.Errorf("failed to get courier: %w", err)
	}

	if courierLatestPositionResponse != nil {
		locationPosition = &LocationPosition{
			Latitude:  courierLatestPositionResponse.Latitude,
			Longitude: courierLatestPositionResponse.Longitude,
		}
	}

	courierResponse := CourierWithLatestPosition{
		FirstName:      courier.FirstName,
		ID:             courier.ID,
		IsAvailable:    courier.IsAvailable,
		LatestPosition: locationPosition,
	}

	return &courierResponse, nil
}

// NewCourierService creates new courier service
func NewCourierService(courierClient CourierClientInterface, repo CourierRepositoryInterface) *CourierService {
	return &CourierService{
		courierClient:     courierClient,
		courierRepository: repo,
	}
}
