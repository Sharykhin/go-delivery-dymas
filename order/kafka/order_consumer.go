package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

type OrderConsumer struct {
	orderRepository domain.OrderRepository
}

// NewOrderConsumer create order consumer
func NewOrderConsumer(
	orderRepository domain.OrderRepository,
) *OrderConsumer {
	orderConsumer := &OrderConsumer{
		orderRepository: orderRepository,
	}

	return orderConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumer *OrderConsumer) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var order domain.Order

	if err := json.Unmarshal(message.Value, &order); err != nil {
		log.Printf("failed to unmarshal Kafka message into order struct: %v\n", err)

		return nil
	}

	err, _ := orderConsumer.orderRepository.ApplyCourierToOrder(ctx, &order)

	if err != nil {
		return fmt.Errorf("failed to save a order in the repository: %w", err)
	}

	return nil
}
