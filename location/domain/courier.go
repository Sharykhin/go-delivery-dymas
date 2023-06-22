package domain

import (
	"context"
	"time"
)

type CourierServiceInterface interface {
	SendData(data *CourierLocationEvent) error
}
type CourierLocationEvent struct {
	CourierID string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type CourierRepositoryInterface interface {
	SaveLatestCourierGeoPosition(ctx context.Context, data *CourierLocationEvent) error
}
