package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

const cgroup = "latest_position_courier"

type CourierLocationMessageJsonHandler struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
	brokers string,
	verbose bool,
	oldest bool,
	assignor string,
) (*pkgkafka.Consumer, error) {
	courierLocationMessageJsonHandler := &CourierLocationMessageJsonHandler{
		courierLocationRepository: courierLocationRepository,
	}
	return pkgkafka.NewConsumer(
		courierLocationMessageJsonHandler,
		brokers,
		verbose,
		oldest,
		assignor,
		'latest_position_courier',
	)
}

func (handlerMessage *CourierLocationMessageJsonHandler) HandleJsonMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var courierLocation domain.CourierLocation
	if err := json.Unmarshal(message.Value, &courierLocation); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message into courier location struct: %w", err)
	}
	err := handlerMessage.courierLocationRepository.SaveLatestCourierGeoPosition(ctx, &courierLocation)
	if err != nil {
		return fmt.Errorf("failed to save a courier location in the repository: %w", err)
	}
	return nil
}
