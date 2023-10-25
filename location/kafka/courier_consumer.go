package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	kafkapkg "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
	"log"
)

const cgroup = "latest_position_courier"

type CourierLocationConsumer struct {
	consumerGroup             kafkapkg.ConsumerGroup
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
		kafkapkg.withBrokers(brokers),
		kafkapkg.withVerbose(verbose),
		kafkapkg.withOldest(oldest),
		kafkapkg.withAssignor(assignor),
	)
	return &CourierLocationConsumer{
		consumerGroup:             consumerGroup,
		courierLocationRepository: courierLocationRepository,
	}, err
}

func (courierLocationConsumer *CourierLocationConsumer) ConsumeCourierLatestCourierGeoPositionMessage(ctx context.Context) error {
	return courierLocationConsumer.consumerGroup.ConsumeMessage(ctx, courierLocationConsumer)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (courierLocationConsumer *CourierLocationConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(courierLocationConsumer.consumerGroup.Ready)
	return nil
}

func (courierLocationConsumer *CourierLocationConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (courierLocationConsumer *CourierLocationConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			var courierLocation domain.CourierLocation
			if err := json.Unmarshal(message.Value, &courierLocation); err != nil {
				log.Printf("Failed to unmarshal Kafka message into courier location struct: %v\n", err)
				session.MarkMessage(message, "")
				return nil
			}
			err := courierLocationConsumer.courierLocationRepository.SaveLatestCourierGeoPosition(session.Context(), &courierLocation)
			if err != nil {
				log.Printf("Failed to save a courier location in the repository: %v", err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
