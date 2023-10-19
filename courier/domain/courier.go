package domain

import (
	"context"
	"errors"
	"fmt"
)

var ErrCourierNotFound = errors.New("courier was not found")

type CourierWithLatestPosition struct {
	LatestPosition *LocationPosition
	FirstName      string
	ID             string
	IsAvailable    bool
}
type CourierClientInterface interface {
	GetLatestPosition(ctx context.Context, courierID string) (*LocationPosition, error)
}
type LocationPosition struct {
	Latitude  float64
	Longitude float64
}
type CourierService struct {
	courierClient     CourierClientInterface
	courierRepository CourierRepositoryInterface
}

type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
}

type CourierServiceInterface interface {
	GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierWithLatestPosition, error)
}

type Courier struct {
	ID          string
	FirstName   string
	IsAvailable bool
}

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

func NewCourierService(CourierClient CourierClientInterface, repo CourierRepositoryInterface) *CourierService {
	return &CourierService{
		courierClient:     CourierClient,
		courierRepository: repo,
	}
}
