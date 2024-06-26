package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationConsumer struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
) *CourierLocationConsumer {
	courierLocationConsumer := &CourierLocationConsumer{
		courierLocationRepository: courierLocationRepository,
	}

	return courierLocationConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (courierLocationConsumer *CourierLocationConsumer) HandleJSONMessage(ctx context.Context, message []byte) error {
	latestCourierLocationMessage := avro.NewLatestCourierLocationMessage()
	var courierLocation domain.CourierLocation
	if err := latestCourierLocationMessage.UnmarshalJSON(message); err != nil {
		log.Printf("failed to unmarshal Kafka message into courier location struct: %v\n", err)

		return nil
	}

	time := time.UnixMilli(latestCourierLocationMessage.Created_at)
	courierLocation = domain.CourierLocation{
		CourierID: latestCourierLocationMessage.Courier_id,
		Latitude:  latestCourierLocationMessage.Latitude,
		Longitude: latestCourierLocationMessage.Longitude,
		CreatedAt: time,
	}

	err := courierLocationConsumer.courierLocationRepository.SaveLatestCourierGeoPosition(ctx, &courierLocation)

	if err != nil {
		return fmt.Errorf("failed to save a courier location in the repository: %w", err)
	}

	return nil
}
