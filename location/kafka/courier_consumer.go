package kafka

import (
	"context"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationConsumer struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
) *CourierLocationConsumer {

	return &CourierLocationConsumer{
		courierLocationRepository: courierLocationRepository,
	}
}

func (c *CourierLocationConsumer) HandleMessage(ctx context.Context, payload domain.CourierLocation) error {
	// Do your stuff

	return nil
}
