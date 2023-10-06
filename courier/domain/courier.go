package domain

import (
	"context"
)

type LocationPosition struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}
type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
}

type LocationPositionServiceInterface interface {
	GetCourierLatestPosition(ctx context.Context, courierID string) (*LocationPosition, error)
}

type Courier struct {
	Id          string
	FirstName   string
	IsAvailable bool
}
