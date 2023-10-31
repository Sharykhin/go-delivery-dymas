package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationMessageJsonHandler struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
) *CourierLocationMessageJsonHandler {
	courierLocationMessageJsonHandler := &CourierLocationMessageJsonHandler{
		courierLocationRepository: courierLocationRepository,
	}

	return courierLocationMessageJsonHandler
}

func (handlerMessage *CourierLocationMessageJsonHandler) HandleJsonMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var courierLocation domain.CourierLocation
	if err := json.Unmarshal(message.Value, &courierLocation); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message into courier location struct: %w", err)
	}
	err := handlerMessage.courierLocationRepository.SaveLatestCourierGeoPosition(ctx, &courierLocation)
	if err != nil {
		return fmt.Errorf("failed to save a courier location in the repository: %w", err)
	}
	return nil
}
