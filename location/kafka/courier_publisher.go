package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

type CourierLocationLatestPublisher struct {
	publisher *pkgkafka.Publisher
}

func NewCourierLocationPublisher(publisher *pkgkafka.Publisher) *CourierLocationLatestPublisher {
	courierPublisher := CourierLocationLatestPublisher{}
	courierPublisher.publisher = publisher
	return &courierPublisher
}

func (courierPublisher *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {
	message, err := json.Marshal(courierLocation)

	if err != nil {
		return fmt.Errorf("failed to marshal courier location before sending Kafka event: %w", err)
	}

	if courierPublisher.publisher != nil {
		err = courierPublisher.publisher.PublishMessage(ctx, message)

		if err != nil {
			return fmt.Errorf("failed to publish courier location: %w", err)
		}
	} else {
		return errors.New("courier publisher was not found")
	}

	return nil
}
