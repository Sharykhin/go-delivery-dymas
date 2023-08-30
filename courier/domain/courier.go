package domain

import (
	"context"
)

type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier Courier) (Courier, error)
}

type Courier struct {
	Id          string `json:"id" validate:"uuid"`
	FirstName   string `json:"first_name" validate:"required"`
	IsAvailable bool   `json:"is_available" validate:"boolean"`
}
