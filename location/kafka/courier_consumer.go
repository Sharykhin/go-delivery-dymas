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
	latestCourierLocation     *avro.LatestCourierLocation
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
	latestCourierLocation *avro.LatestCourierLocation,
) *CourierLocationConsumer {
	courierLocationConsumer := &CourierLocationConsumer{
		courierLocationRepository: courierLocationRepository,
		latestCourierLocation:     latestCourierLocation,
	}

	return courierLocationConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (courierLocationConsumer *CourierLocationConsumer) HandleJSONMessage(ctx context.Context, message []byte) error {
	var courierLocation domain.CourierLocation
	if err := courierLocationConsumer.latestCourierLocation.UnmarshalJSON(message); err != nil {
		log.Printf("failed to unmarshal Kafka message into courier location struct: %v\n", err)

		return nil
	}

	time := time.UnixMilli(courierLocationConsumer.latestCourierLocation.Created_at)
	courierLocation = domain.CourierLocation{
		CourierID: courierLocationConsumer.latestCourierLocation.Courier_id,
		Latitude:  courierLocationConsumer.latestCourierLocation.Latitude,
		Longitude: courierLocationConsumer.latestCourierLocation.Longitude,
		CreatedAt: time,
	}

	err := courierLocationConsumer.courierLocationRepository.SaveLatestCourierGeoPosition(ctx, &courierLocation)

	if err != nil {
		return fmt.Errorf("failed to save a courier location in the repository: %w", err)
	}

	return nil
}
