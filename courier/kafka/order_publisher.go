package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

const OrderTopicValidation = "order_validations"

// OrderValidationPublisher publisher for kafka
type OrderValidationPublisher struct {
	publisher *pkgkafka.Publisher
}

// NewOrderValidationPublisher creates new publisher and init
func NewOrderValidationPublisher(publisher *pkgkafka.Publisher) *OrderValidationPublisher {
	orderValidationPublisher := OrderValidationPublisher{
		publisher: publisher,
	}

	return &orderValidationPublisher
}

// PublishValidationResult sends order message in json format in Kafka.
func (orderPublisher *OrderValidationPublisher) PublishValidationResult(ctx context.Context, courierAssigment *domain.CourierAssignments) error {
	messageOrder := domain.OrderMessageValidation{
		IsSuccessful: true,
		ServiceName:  "courier",
		Payload: domain.CourierAssignments{
			OrderID:   courierAssigment.OrderID,
			CourierID: courierAssigment.CourierID,
		},
	}

	message, err := json.Marshal(messageOrder)

	if err != nil {
		return fmt.Errorf("failed to marshal order before sending Kafka event: %w", err)
	}

	err = orderPublisher.publisher.PublishMessage(ctx, message, []byte(courierAssigment.OrderID))

	if err != nil {
		return fmt.Errorf("failed to publish order event: %w", err)
	}

	return nil
}
