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
	orderClient     domain.OrderClient
}

// NewOrderConsumer create order consumer
func NewOrderConsumer(
	orderRepository domain.OrderRepository,
	orderPublisher domain.OrderPublisher,
	orderClient domain.OrderClient,
) *OrderConsumer {
	orderConsumer := &OrderConsumer{
		orderRepository: orderRepository,
		orderPublisher:  orderPublisher,
		orderClient:     orderClient,
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

	order, err := orderConsumer.orderClient.GetAssignCourier(ctx, orderMessage.Payload)

	if err != nil {
		return fmt.Errorf("failed to get assign courier: %w", err)
	}

	if orderMessage.Event == domain.MessageStatusUpdated {
		return nil
	}

	_, err = orderConsumer.orderRepository.AssignCourierToOrder(ctx, order)

	if err != nil {
		return fmt.Errorf("failed to save a order in the repository: %w", err)
	}

	orderConsumer.orderPublisher.PublishOrder(ctx, orderMessage.Payload, domain.MessageStatusUpdated)

	return nil
}
