package domain

import (
	"context"
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
