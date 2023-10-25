package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type JSONHandler interface {
	HandleJSONMessage(ctx context.Context, payload any) error
}

type Consumer struct {
	verboseEnabled bool
	assignor       string
	topic          string
	keepRunning    bool
	ready          chan bool

	consumerGroup sarama.ConsumerGroup
	jsonHandler   JSONHandler
}

func WithVerboseConsumer(isEnabled bool) func(*Consumer) {
	return func(consumer *Consumer) {
		consumer.verboseEnabled = isEnabled
	}
}

func NewConsumer(topic string, opts ...func(*Consumer)) *Consumer {
	consumer := &Consumer{
		topic:       topic,
		keepRunning: true,
		ready:       make(chan bool),
	}
	// create consumer group
	for _, opt := range opts {
		opt(consumer)
	}

	if consumer.verboseEnabled {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	switch consumer.assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		err := fmt.Errorf("unrecognized consumer group partition assignor: %s", consumer.assignor)

		log.Panic(err)
	}

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split("", ","), consumer.topic, config)
	if err != nil {
		log.Panic(err)
	}

	consumer.consumerGroup = consumerGroup

	return consumer
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)

	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) RegisterJSONHandler(ctx context.Context, handler JSONHandler) {
	c.jsonHandler = handler
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
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
			if err := c.consumerGroup.Consume(ctx, strings.Split(c.topic, ","), c); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Panicf("Error from consumer: %v", ctx.Err())
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for c.keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			c.keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			c.keepRunning = false
		case <-sigusr1:
			c.toggleConsumptionFlow(c.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		return fmt.Errorf("Error closing client: %w", err)
	}

	return nil
}

func (c *Consumer) toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			err := c.jsonHandler.HandleJSONMessage(session.Context(), message.Value)
			if err != nil {
				log.Printf("Failed to save a courier location in the repository: %v", err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
