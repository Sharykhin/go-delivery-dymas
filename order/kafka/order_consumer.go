package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderConsumerValidation consumes message order validation from kafka
type OrderConsumerValidation struct {
	orderService domain.OrderService
}

// CourierPayload imagines contract how view courier payload from third system
type CourierPayload struct {
	CourierID          string    `json:"courier_id"`
	CourierValidatedAt time.Time `json:"courier_validated_at"`
}

// OrderMessageValidation sends in third system for service information about order assign.
type OrderMessageValidation struct {
	IsSuccessful bool            `json:"isSuccessful"`
	Payload      json.RawMessage `json:"payload"`
	ServiceName  string          `json:"serviceName"`
	OrderID      string          `json:"order_id"`
}

// NewOrderConsumerValidation creates order validation consumer
func NewOrderConsumerValidation(
	orderService domain.OrderService,
) *OrderConsumerValidation {
	orderConsumer := &OrderConsumerValidation{
		orderService: orderService,
	}

	return orderConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumerValidation *OrderConsumerValidation) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var orderMessageValidation OrderMessageValidation
	if err := json.Unmarshal(message.Value, &orderMessageValidation); err != nil {
		log.Printf("failed to unmarshal Kafka message into order validation struct: %v\n", err)

		return nil
	}

	err := orderConsumerValidation.orderService.ValidateOrderForService(
		ctx,
		orderMessageValidation.ServiceName,
		orderMessageValidation.OrderID,
		orderMessageValidation.Payload,
	)

	if err != nil {
		return fmt.Errorf("failed to validate order: %w", err)
	}

	return nil
}
