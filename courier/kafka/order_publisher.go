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

// OrderMessageValidation sends in third system for service information about order assign.
type OrderMessageValidation struct {
	IsSuccessful bool                     `json:"isSuccessful"`
	Payload      domain.CourierAssignment `json:"payload"`
	ServiceName  string                   `json:"serviceName"`
	event        string                   `json:"event"`
}

// NewOrderValidationPublisher creates new publisher and init
func NewOrderValidationPublisher(publisher *pkgkafka.Publisher) *OrderValidationPublisher {
	orderValidationPublisher := OrderValidationPublisher{
		publisher: publisher,
	}

	return &orderValidationPublisher
}

// PublishValidationResult sends order message in json format in Kafka.
func (orderPublisher *OrderValidationPublisher) PublishValidationResult(ctx context.Context, courierAssigment *domain.CourierAssignment) error {
	messageOrderValidation := OrderMessageValidation{
		IsSuccessful: true,
		ServiceName:  "courier",
		Payload: domain.CourierAssignment{
			OrderID:   courierAssigment.OrderID,
			CourierID: courierAssigment.CourierID,
		},
	}

	message, err := json.Marshal(messageOrderValidation)

	if err != nil {
		return fmt.Errorf("failed to marshal order message validation before sending Kafka event: %w", err)
	}

	err = orderPublisher.publisher.PublishMessage(ctx, message, []byte(courierAssigment.OrderID))

	if err != nil {
		return fmt.Errorf("failed to publish order event: %w", err)
	}

	return nil
}
