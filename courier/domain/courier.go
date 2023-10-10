package domain

import (
	"context"
)

type CourierResponse struct {
	LatestPosition *LocationPosition `json:"last_position"`
	FirstName      string            `json:"first_name" validate:"required"`
	Id             string            `json:"id" validate:"uuid,required"`
	IsAvailable    bool              `json:"is_available" validate:"boolean,required"`
}
type CourierLocationPositionClientInterface interface {
	GetCourierLatestPosition(ctx context.Context, courierID string) (*LocationPosition, error)
}
type LocationPosition struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}
type LocationPositionService struct {
	CourierClient     CourierLocationPositionClientInterface
	courierRepository CourierRepositoryInterface
}

type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
}

type LocationPositionServiceInterface interface {
	GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierResponse, error)
}

type Courier struct {
	Id          string
	FirstName   string
	IsAvailable bool
}

func (s LocationPositionService) GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierResponse, error) {
	courier, err := s.courierRepository.GetCourierByID(
		ctx,
		courierID,
	)
	if err != nil {
		return nil, err
	}
	courierLatestPositionResponse, err := s.CourierClient.GetCourierLatestPosition(ctx, courierID)
	if err != nil {
		return nil, err
	}
	locationPosition := LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}
	courierResponse := CourierResponse{
		FirstName:      courier.FirstName,
		Id:             courier.Id,
		IsAvailable:    courier.IsAvailable,
		LatestPosition: &locationPosition,
	}
	return &courierResponse, nil
}

func NewLocationPositionService(CourierClient CourierLocationPositionClientInterface) *LocationPositionService {
	return &LocationPositionService{
		CourierClient: CourierClient,
	}
}
