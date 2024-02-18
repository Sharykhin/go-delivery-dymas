package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderConsumerValidation consumes message order validation from kafka
type OrderConsumerValidation struct {
	orderRepository domain.OrderRepository
	orderPublisher  domain.OrderPublisher
}

// OrderMessageValidation sends in third system for service information about order assign.
type OrderMessageValidation struct {
	IsSuccessful bool            `json:"isSuccessful"`
	Payload      json.RawMessage `json:"payload"`
	ServiceName  string          `json:"serviceName"`
}

// NewOrderConsumerValidation creates order validation consumer
func NewOrderConsumerValidation(
	orderRepository domain.OrderRepository,
	orderPublisher domain.OrderPublisher,
) *OrderConsumerValidation {
	orderConsumer := &OrderConsumerValidation{
		orderRepository: orderRepository,
		orderPublisher:  orderPublisher,
	}

	return orderConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumerValidation *OrderConsumerValidation) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var orderMessageValidation OrderMessageValidation

	if err := json.Unmarshal(message.Value, &orderMessageValidation); err != nil {
		log.Printf("failed to unmarshal Kafka message into order struct: %v\n", err)

		return nil
	}

	switch orderMessageValidation.ServiceName {
	case "courier":
		var courierPayload domain.CourierPayload
		json.Unmarshal(orderMessageValidation.Payload, &courierPayload)
		order, err := orderConsumerValidation.orderRepository.ChangeOrderStatusAfterValidation(
			ctx,
			&courierPayload,
			orderMessageValidation.IsSuccessful,
			domain.OrderStatusAccepted,
			"courier",
		)

		if err != nil {
			return fmt.Errorf("failed to save a order in the repository: %w", err)
		}

		return orderConsumerValidation.publishOrderMessage(ctx, order)
	}

	return nil
}

func (orderConsumerValidation *OrderConsumerValidation) publishOrderMessage(ctx context.Context, order domain.Order) error {
	err := orderConsumerValidation.orderPublisher.PublishOrder(ctx, &order, domain.EventOrderUpdated)

	if err != nil {
		return fmt.Errorf("failed to publish a order in the kafka: %w", err)
	}

	return nil
}
