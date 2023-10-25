package kafka

import (
	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	kafkapkg "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
	"log"
)

const cgroup = "latest_position_courier"

type CourierLocationConsumer struct {
	kafkapkg.ConsumerGroup
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
	brokers string,
	verbose bool,
	oldest bool,
	assignor string,
) (*CourierLocationConsumer, error) {

	consumerGroup, err := kafkapkg.NewGroupConsumer(
		kafkapkg.CourierLocationConsumer{},
		kafkapkg.WithBrokers(brokers),
		kafkapkg.WithVerbose(verbose),
		kafkapkg.WithOldest(oldest),
		kafkapkg.WithAssignor(assignor),
	)
	return &CourierLocationConsumer{
		consumerGroup:             consumerGroup,
		courierLocationRepository: courierLocationRepository,
	}, err
}

func (courierLocationConsumer *CourierLocationConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	return courierLocationConsumer.HandleJsonMessage(ctx, claim, domain.CourierLocation{}, func() error {
		err := courierLocationConsumer.courierLocationRepository.SaveLatestCourierGeoPosition(session.Context(), &courierLocation)
		if err != nil {
			log.Printf("Failed to save a courier location in the repository: %v", err)
		}

		return err
	})
}
