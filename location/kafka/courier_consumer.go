package kafka

import (
	"context"
	"encoding/json"
	
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

func (c *CourierLocationConsumer) HandleJSONMessage(ctx context.Context, payload []byte) error {
	var location domain.CourierLocation
	err := json.Unmarshal(payload, &location)
	if err != nil {
		return nil
	}

	return nil
}
