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
	courierService domain.CourierService
}

// OrderPayload  needs for order message
type OrderPayload struct {
	OrderID string `json:"id"`
}

// OrderMessage will consume, when order create and publish in queue.
type OrderMessage struct {
	OrderPayload OrderPayload `json:"payload"`
	Event        string       `json:"event"`
}

// NewOrderConsumer creates and init order consumer this consumer consume message from kafka
func NewOrderConsumer(
	courierService domain.CourierService,
) *OrderConsumer {
	courierConsumer := &OrderConsumer{
		courierService: courierService,
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

	err := orderConsumer.courierService.AssignOrderToCourier(ctx, orderMessage.OrderPayload.OrderID)
	if err != nil {
		return fmt.Errorf("can not assign order to courier: %w", err)
	}

	return nil
}
