package domain

import (
	"context"
	"time"
)

type CourierServiceInterface interface {
	SendData(data *CourierRepositoryData, ctx context.Context) error
}
type MessageKafka struct {
	CourierId string    `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}
type CourierRepositoryData struct {
	CourierID string
	Latitude  float64
	Longitude float64
}

type CourierRepositoryInterface interface {
	SaveLatestCourierGeoPosition(data *CourierRepositoryData, ctx context.Context) error
}
