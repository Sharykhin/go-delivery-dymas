package domain

import (
	"context"
)

type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierById(ctx context.Context, courierID string) (*Courier, error)
}

type Courier struct {
	Id          string
	FirstName   string
	IsAvailable bool
}
