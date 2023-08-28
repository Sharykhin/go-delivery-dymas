package domain

import (
	"context"
	"fmt"
	"time"
)

type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier CourierModel) error
}

type CourierModel struct {
	Id          string `json:"id" validate:"required,uuid"`
	FirstName   string `json:"first_name" validate:"required"`
	IsAvailable bool   `json:"is_available" validate:"required,boolean"`
}
