package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationLatestPublisher struct {
	publisher Publisher
}

func NewCourierLocationPublisher(publisher Publisher) *CourierLocationLatestPublisher {
	courierPublisher := CourierLocationLatestPublisher{
		publisher: publisher,
	}

	return &courierPublisher
}

func (p *CourierLocationLatestPublisher) PublishLatestCourierLocation(ctx context.Context, courierLocation *domain.CourierLocation) error {

	err := p.publisher.SendJSONMessage(ctx, courierLocation)
	if err != nil {
		return fmt.Errorf("faield to send kafka event as json: %w", err)
	}

	return nil
}
