package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderConsumerValidation consumes message order validation from kafka
type OrderConsumerValidation struct {
	orderService           domain.OrderService
	orderValidationMessage avro.OrderValidationMessage
}

// OrderMessageValidation sends in third system for service information about order assign.
type OrderMessageValidation struct {
	IsSuccessful bool            `json:"is_successful"`
	Payload      json.RawMessage `json:"payload"`
	ServiceName  string          `json:"service_name"`
	OrderID      string          `json:"order_id"`
}

// NewOrderConsumerValidation creates order validation consumer
func NewOrderConsumerValidation(
	orderService domain.OrderService,
	orderValidationMessage avro.OrderValidationMessage,
) *OrderConsumerValidation {
	orderConsumer := &OrderConsumerValidation{
		orderService:           orderService,
		orderValidationMessage: orderValidationMessage,
	}

	return orderConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumerValidation *OrderConsumerValidation) HandleJSONMessage(ctx context.Context, message []byte) error {
	if err := orderConsumerValidation.orderValidationMessage.UnmarshalJSON(message); err != nil {
		log.Printf("failed to unmarshal Kafka message into order validation struct: %v\n", err)

		return nil
	}

	payload := domain.OrderValidationPayload{
		CourierID: orderConsumerValidation.orderValidationMessage.Payload.Courier_id.String,
	}
	err := orderConsumerValidation.orderService.ValidateOrderForService(
		ctx,
		orderConsumerValidation.orderValidationMessage.Service_name,
		orderConsumerValidation.orderValidationMessage.Order_id,
		payload,
	)

	if err != nil {
		return fmt.Errorf("failed to validate order: %w", err)
	}

	return nil
}
