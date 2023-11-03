package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationConsumer struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
	errorHandler              domain.TypeErrorHandler
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
	h domain.TypeErrorHandler,
) *CourierLocationConsumer {
	courierLocationConsumer := &CourierLocationConsumer{
		courierLocationRepository: courierLocationRepository,
		errorHandler:              h,
	}

	return courierLocationConsumer
}

func (courierLocationConsumer *CourierLocationConsumer) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var courierLocation domain.CourierLocation

	if err := json.Unmarshal(message.Value, &courierLocation); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message into courier location struct: %w", err)
	}

	err := courierLocationConsumer.courierLocationRepository.SaveLatestCourierGeoPosition(ctx, &courierLocation)

	if err != nil {
		return fmt.Errorf("failed to save a courier location in the repository: %w", err)
	}

	return nil
}

func (courierLocationConsumer *CourierLocationConsumer) IsRetryAttempt(err error) bool {
	return courierLocationConsumer.errorHandler.CompareTypeError(err)
}
