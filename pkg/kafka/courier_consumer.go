package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type CourierLocationConsumer struct {
	consumerGroup             sarama.ConsumerGroup
	keepRunning               bool
	courierLocationRepository domain.CourierLocationRepositoryInterface
	ready                     chan bool
	topic                     string
}

func NewCourierLocationConsumer(
	courierLocationRepository domain.CourierLocationRepositoryInterface,
	brokers string,
	verbose bool,
	oldest bool,
	assignor string,
) (*CourierLocationConsumer, error) {

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		err := fmt.Errorf("Unrecognized consumer group partition assignor: %s", assignor)

		return nil, err
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), cgroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create courier location consumer: %w", err)
	}
	return &CourierLocationConsumer{
		consumerGroup:             consumerGroup,
		keepRunning:               true,
		ready:                     make(chan bool),
		courierLocationRepository: courierLocationRepository,
	}, nil
}

func (courierLocationConsumer *CourierLocationConsumer) ConsumeCourierLatestCourierGeoPositionMessage(ctx context.Context) error {
	consumptionIsPaused := false
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := courierLocationConsumer.consumerGroup.Consume(ctx, strings.Split(CourierLocationConsumer.topic, ","), courierLocationConsumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Panicf("Error from consumer: %v", ctx.Err())
				return
			}
			courierLocationConsumer.ready = make(chan bool)
		}
	}()

	<-courierLocationConsumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for courierLocationConsumer.keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			courierLocationConsumer.keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			courierLocationConsumer.keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(courierLocationConsumer.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := courierLocationConsumer.consumerGroup.Close(); err != nil {
		return fmt.Errorf("Error closing client: %w", err)
	}

	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (courierLocationConsumer *CourierLocationConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(courierLocationConsumer.ready)
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

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}