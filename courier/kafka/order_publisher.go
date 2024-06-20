package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

const OrderTopicValidation = "order_validations.v1"

// OrderValidationPublisher publisher for kafka
type OrderValidationPublisher struct {
	publisher *pkgkafka.Publisher
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
func NewOrderValidationPublisher(publisher *pkgkafka.Publisher) *OrderValidationPublisher {
	orderValidationPublisher := OrderValidationPublisher{
		publisher: publisher,
	}

	return &orderValidationPublisher
}

// PublishValidationResult sends order message in json format in Kafka.
func (orderPublisher *OrderValidationPublisher) PublishValidationResult(ctx context.Context, courierAssigment *domain.CourierAssignment) error {
	OrderValidationMessage := avro.NewOrderValidationMessage()
	OrderValidationMessage.Order_id = courierAssigment.OrderID
	OrderValidationMessage.Service_name = "courier"
	OrderValidationMessage.Is_successful = true
	OrderValidationMessage.Payload.Courier_id.String = courierAssigment.CourierID
	OrderValidationMessage.Payload.Courier_id.Null = nil

	message, err := OrderValidationMessage.MarshalJSON()

	if err != nil {
		return fmt.Errorf("failed to marshal order message validation before sending Kafka event: %w", err)
	}

	schema := OrderValidationMessage.Schema()
	err = orderPublisher.publisher.PublishMessage(ctx, message, []byte(courierAssigment.OrderID), schema)

	if err != nil {
		return fmt.Errorf("failed to publish order message validation event: %w", err)
	}

	return nil
}
