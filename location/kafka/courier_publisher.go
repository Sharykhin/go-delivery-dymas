package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

// CourierLocationLatestPublisher Publisher for kafka
type CourierLocationLatestPublisher struct {
	publisher             *pkgkafka.Publisher
	latestCourierLocation *avro.LatestCourierLocation
}

// NewCourierLocationPublisher creates new publisher and init
func NewCourierLocationPublisher(publisher *pkgkafka.Publisher, latestCourierLocation *avro.LatestCourierLocation) *CourierLocationLatestPublisher {
	courierPublisher := CourierLocationLatestPublisher{}
	courierPublisher.publisher = publisher
	courierPublisher.latestCourierLocation = latestCourierLocation
	return &courierPublisher
}

// PublishLatestCourierLocation sends latest courier position message in json format in Kafka.
func (courierPublisher *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	courierPublisher.latestCourierLocation.Courier_id = courierLocation.CourierID
	courierPublisher.latestCourierLocation.Longitude = courierLocation.Longitude
	courierPublisher.latestCourierLocation.Latitude = courierLocation.Latitude
	courierPublisher.latestCourierLocation.Created_at = courierLocation.CreatedAt.Unix()
	message, err := courierPublisher.latestCourierLocation.MarshalJSON()

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}

	schema := courierPublisher.latestCourierLocation.Schema()
	err = courierPublisher.publisher.PublishMessage(ctx, message, []byte(courierLocation.CourierID), schema)

	if err != nil {
		return fmt.Errorf("failed to publish courier location: %w", err)
	}

	return nil
}
