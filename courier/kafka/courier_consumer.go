package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
)

// OrderTopic where we have message with different event for order
const OrderTopic = "orders"

// OrderConsumer gets order from kafka and apply order to courier and send order message validations
type OrderConsumer struct {
	courierRepository        domain.CourierRepositoryInterface
	orderValidationPublisher domain.OrderValidationPublisher
}

// Customer needs for order message
type Customer struct {
	PhoneNumber string `json:"phone_number"`
}

// OrderPayload  needs for order message
type OrderPayload struct {
	OrderID  string   `json:"id"`
	Customer Customer `json:"customer"`
}

// OrderMessage will consume, when order create and publish in queue.
type OrderMessage struct {
	OrderPayload OrderPayload `json:"payload"`
	Event        string       `json:"event"`
}

// NewOrderConsumer creates and init order consumer this consumer consume message from kafka
func NewOrderConsumer(
	courierRepository domain.CourierRepositoryInterface,
	orderValidationPublisher domain.OrderValidationPublisher,
) *OrderConsumer {
	courierConsumer := &OrderConsumer{
		courierRepository:        courierRepository,
		orderValidationPublisher: orderValidationPublisher,
	}

	return courierConsumer
}

// HandleJSONMessage Handle kafka message in json format
func (orderConsumer *OrderConsumer) HandleJSONMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var orderMessage OrderMessage
	if err := json.Unmarshal(message.Value, &orderMessage); err != nil {
		log.Printf("failed to unmarshal Kafka message into courier order message struct: %v\n", err)

		return nil
	}

	if orderMessage.Event == "updated" {
		return nil
	}

	courierAssigment, err := orderConsumer.courierRepository.AssignOrderToCourier(ctx, orderMessage.OrderPayload.OrderID)
	if err != nil {
		return fmt.Errorf("failed to save a courier assigments in the repository: %w", err)
	}

	err = orderConsumer.orderValidationPublisher.PublishValidationResult(ctx, &courierAssigment)

	if err != nil {
		return fmt.Errorf("failed to publish a order message validation in kafka: %w", err)
	}

	return nil
}
