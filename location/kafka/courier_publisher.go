package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

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
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}

	err = courierPublisher.publisher.PublishMessage(ctx, message, []byte(courierLocation.CourierID))

	if err != nil {
		return fmt.Errorf("failed to publish courier location: %w", err)
	}

	return nil
}
