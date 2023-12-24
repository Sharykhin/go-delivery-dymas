package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

// OrderPublisher publisher for kafka
type OrderPublisher struct {
	publisher *pkgkafka.Publisher
}

// NewOrderPublisher creates new publisher and init
func NewOrderPublisher(publisher *pkgkafka.Publisher) *OrderPublisher {
	orderPublisher := OrderPublisher{
		publisher: publisher,
	}

	return &orderPublisher
}

// PublishOrder sends order message in json format in Kafka.
func (orderPublisher *OrderPublisher) PublishOrder(ctx context.Context, order *domain.Order) error {
	message, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to marshal order before sending Kafka event: %w", err)
	}

	err = orderPublisher.publisher.PublishMessage(ctx, message, order.ID)

	if err != nil {
		return fmt.Errorf("failed to publish order event: %w", err)
	}

	return nil
}
