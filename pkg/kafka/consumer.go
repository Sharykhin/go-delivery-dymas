package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type ConsumerGroupInterface interface {
}

type ConsumerGroup struct {
	verbose       bool
	assignor      string
	oldest        bool
	brokers       string
	keepRunning   bool
	consumerGroup sarama.ConsumerGroup
	Ready         chan bool
}

func NewGroupConsumer(cgroup string, opts ...func(cg *ConsumerGroup)) (*ConsumerGroup, error) {
	var consumerGroup ConsumerGroup
	for _, opt := range opts {
		opt(&consumerGroup)
	}
	if consumerGroup.verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	switch consumerGroup.assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		err := fmt.Errorf("Unrecognized consumer group partition assignor: %s", consumerGroup.assignor)

		return nil, err
	}

	if consumerGroup.oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	cg, err := sarama.NewConsumerGroup(strings.Split(consumerGroup.brokers, ","), cgroup, config)
	consumerGroup.consumerGroup = cg
	return &consumerGroup, err
}

func withAssignor(assignor string) func(cg *ConsumerGroup) {
	return func(cg *ConsumerGroup) {
		cg.assignor = assignor
	}
}

func withVerbose(verbose bool) func(cg *ConsumerGroup) {
	return func(cg *ConsumerGroup) {
		cg.verbose = verbose
	}
}

func withOldest(oldest bool) func(cg *ConsumerGroup) {
	return func(cg *ConsumerGroup) {
		cg.oldest = oldest
	}
}

func withBrokers(brokers string) func(cg *ConsumerGroup) {
	return func(cg *ConsumerGroup) {
		cg.brokers = brokers
	}
}

func (consumerGroup *ConsumerGroup) ConsumeMessage(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
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
			if err := consumerGroup.consumerGroup.Consume(ctx, strings.Split(topic, ","), handler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Panicf("Error from consumer: %v", ctx.Err())
				return
			}
			consumerGroup.Ready = make(chan bool)
		}
	}()

	<-consumerGroup.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for consumerGroup.keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			consumerGroup.keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			consumerGroup.keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(consumerGroup.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := consumerGroup.consumerGroup.Close(); err != nil {
		return fmt.Errorf("Error closing client: %w", err)
	}

	return nil
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
