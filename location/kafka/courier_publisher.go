package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

const LatestPositionCourierTopic = "latest_position_courier.v1"

// CourierLocationLatestPublisher Publisher for kafka
type CourierLocationLatestPublisher struct {
	publisher *pkgkafka.Publisher
}

// NewCourierLocationPublisher creates new publisher and init
func NewCourierLocationPublisher(publisher *pkgkafka.Publisher) *CourierLocationLatestPublisher {
	courierPublisher := CourierLocationLatestPublisher{}
	courierPublisher.publisher = publisher
	return &courierPublisher
}

// PublishLatestCourierLocation sends latest courier position message in json format in Kafka.
func (courierPublisher *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	latestCourierLocation := avro.NewLatestCourierLocation()
	latestCourierLocation.Courier_id = courierLocation.CourierID
	latestCourierLocation.Longitude = courierLocation.Longitude
	latestCourierLocation.Latitude = courierLocation.Latitude
	latestCourierLocation.Created_at = courierLocation.CreatedAt.Unix()
	message, err := latestCourierLocation.MarshalJSON()

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}

	schema := latestCourierLocation.Schema()
	err = courierPublisher.publisher.PublishMessage(ctx, message, []byte(courierLocation.CourierID), schema)

	if err != nil {
		return fmt.Errorf("failed to publish courier location: %w", err)
	}

	return nil
}
