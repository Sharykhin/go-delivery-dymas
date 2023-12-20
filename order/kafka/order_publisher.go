package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

// OrderPublisherCreate publisher for kafka
type OrderPublisherCreate struct {
	publisher *pkgkafka.Publisher
}

// NewOrderPublisher creates new publisher and init
func NewOrderPublisher(publisher *pkgkafka.Publisher) *OrderPublisherCreate {
	orderPublisherCreate := OrderPublisherCreate{
		publisher: publisher,
	}

	return &orderPublisherCreate
}

// PublishOrder sends order message in json format in Kafka.
func (orderPublisher *OrderPublisherCreate) PublishOrder(ctx context.Context, order *domain.Order) error {
	message, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to marshal order before sending Kafka event: %w", err)
	}

	err = orderPublisher.publisher.PublishMessage(ctx, message)

	if err != nil {
		return fmt.Errorf("failed to publish courier location: %w", err)
	}

	return nil
}
