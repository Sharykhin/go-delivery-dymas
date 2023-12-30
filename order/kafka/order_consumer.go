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
	orderPublisher  domain.OrderPublisher
}

// NewOrderConsumer create order consumer
func NewOrderConsumer(
	orderRepository domain.OrderRepository,
	orderPublisher domain.OrderPublisher,
) *OrderConsumer {
	orderConsumer := &OrderConsumer{
		orderRepository: orderRepository,
		orderPublisher:  orderPublisher,
	}

	return orderConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumer *OrderConsumer) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var orderMessage domain.OrderMessage

	if err := json.Unmarshal(message.Value, &orderMessage); err != nil {
		log.Printf("failed to unmarshal Kafka message into order struct: %v\n", err)

		return nil
	}

	err, _ := orderConsumer.orderRepository.AssignCourierToOrder(ctx, orderMessage.Payload)

	if err != nil {
		return fmt.Errorf("failed to save a order in the repository: %w", err)
	}

	orderConsumer.orderPublisher.PublishOrder(ctx, orderMessage.Payload, domain.MessageStatusUpdated)

	return nil
}
