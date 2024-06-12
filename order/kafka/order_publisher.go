package kafka

import (
	"context"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/avro/v1"
	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

// OrderPublisher publisher for kafka
type OrderPublisher struct {
	publisher    *pkgkafka.Publisher
	orderMessage *avro.OrderMessage
}

// OrderPayload uses for embedding order id and phone customer
type OrderPayload struct {
	OrderID string `json:"id"`
}

// OrderMessage will publish, when order create.
type OrderMessage struct {
	OrderPayload OrderPayload `json:"payload"`
	Event        string       `json:"event"`
}

// NewOrderPublisher creates new publisher and init
func NewOrderPublisher(publisher *pkgkafka.Publisher, orderMessage *avro.OrderMessage) *OrderPublisher {
	orderPublisher := OrderPublisher{
		publisher:    publisher,
		orderMessage: orderMessage,
	}

	return &orderPublisher
}

// PublishOrder sends order message in json format in Kafka.
func (orderPublisher *OrderPublisher) PublishOrder(ctx context.Context, order *domain.Order, event string) error {
	orderPublisher.orderMessage.Payload.Order_id = order.ID
	orderPublisher.orderMessage.Event = event
	message, err := orderPublisher.orderMessage.MarshalJSON()
	schema := orderPublisher.orderMessage.Schema()

	if err != nil {
		return fmt.Errorf("failed to marshal order before sending Kafka event: %w", err)
	}

	err = orderPublisher.publisher.PublishMessage(ctx, message, []byte(order.ID), schema)

	if err != nil {
		return fmt.Errorf("failed to publish order event: %w", err)
	}

	return nil
}
