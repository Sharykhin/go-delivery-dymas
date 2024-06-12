package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

const OrderTopicValidation = "order_validations"

// OrderValidationPublisher publisher for kafka
type OrderValidationPublisher struct {
	publisher              *pkgkafka.Publisher
	orderValidationMessage *avro.OrderValidationMessage
}

// CourierPayload need for send order message validation in kafka
type CourierPayload struct {
	CourierID string `json:"courier_id"`
}

// OrderMessageValidation sends in third system for service information about order assign.
type OrderMessageValidation struct {
	IsSuccessful bool           `json:"is_successful"`
	Payload      CourierPayload `json:"payload"`
	OrderID      string         `json:"order_id"`
	ServiceName  string         `json:"service_name"`
}

// NewOrderValidationPublisher creates new publisher and init
func NewOrderValidationPublisher(publisher *pkgkafka.Publisher, orderValidationMessage *avro.OrderValidationMessage) *OrderValidationPublisher {
	orderValidationPublisher := OrderValidationPublisher{
		publisher:              publisher,
		orderValidationMessage: orderValidationMessage,
	}

	return &orderValidationPublisher
}

// PublishValidationResult sends order message in json format in Kafka.
func (orderPublisher *OrderValidationPublisher) PublishValidationResult(ctx context.Context, courierAssigment *domain.CourierAssignment) error {
	orderPublisher.orderValidationMessage.Order_id = courierAssigment.OrderID
	orderPublisher.orderValidationMessage.Service_name = "courier"
	orderPublisher.orderValidationMessage.Is_successful = true
	orderPublisher.orderValidationMessage.Payload.Courier_id.String = courierAssigment.CourierID

	message, err := orderPublisher.orderValidationMessage.MarshalJSON()

	if err != nil {
		return fmt.Errorf("failed to marshal order message validation before sending Kafka event: %w", err)
	}

	schema := orderPublisher.orderValidationMessage.Schema()
	err = orderPublisher.publisher.PublishMessage(ctx, message, []byte(courierAssigment.OrderID), schema)

	if err != nil {
		return fmt.Errorf("failed to publish order event: %w", err)
	}

	return nil
}
