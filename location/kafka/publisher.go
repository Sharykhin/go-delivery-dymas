package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// Publisher provides describes general API for event-based message communication.
// At the current moment it only supports simple string encoding as JSON.
type Publisher interface {
	SendJSONMessage(ctx context.Context, message any) error
}

type publisher struct {
	topicName     string
	brokerAddress string
	requiredAcks  sarama.RequiredAcks

	producer sarama.AsyncProducer
}

// NewPublisher creates a new instance of Kafka publisher and returns general Publisher interface.
// As for given parameters it only accepts topic name. Managing producer Partitioner and RequiredAcks
// are not supported yet but will be implemented in latter versions.
func NewPublisher(topicName string, opts ...func(*publisher)) Publisher {
	if topicName == "" {
		log.Panic("topic name is required for creating kafka publisher")
	}

	p := &publisher{
		requiredAcks:  sarama.WaitForLocal,
		brokerAddress: os.Getenv("KAFKA_BROKER_ADDRESS"),
		topicName:     topicName,
	}
	for _, opt := range opts {
		opt(p)
	}

	config := sarama.NewConfig()
	if p.brokerAddress == "" {
		log.Panic("kafka broker address is missing, use WithKafkaAddress option or KAFKA_ADDRESS env variable")
	}

	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = p.requiredAcks

	producer, err := sarama.NewAsyncProducer([]string{p.brokerAddress}, config)
	if err != nil {
		log.Panicf("failed to initialize async producer: %v", err)
	}

	return &publisher{
		producer: producer,
	}
}

// WithBrokersAddress provides Kafka broker address. By default, when NewPublisher is called it take it from
// KAFKA_BROKER_ADDRESS environment variable. But if needed to override this use this option.
func WithBrokersAddress(brokerAddress string) func(*publisher) {
	return func(p *publisher) {
		p.brokerAddress = brokerAddress
	}
}

// WithProducerAck set the producer acknowledgements. By default, when NewPublisher is called it sets sarama.WaitForLocal.
// Use this option if you need to override this value.
func WithProducerAck(acks sarama.RequiredAcks) func(*publisher) {
	return func(p *publisher) {
		p.requiredAcks = acks
	}
}

// SendJSONMessage sends message as JSON. Before sending sarama produce message it marshals the given message
// into slice of bytes. Even though the context is not used in this particular method it is defined in the interface,
// and most probably it can be used by anther implementation of event messaging. More over Sarama may introduce
// context supporting in later versions.
func (p *publisher) SendJSONMessage(_ context.Context, message any) error {
	encoded, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal json message before sending Kafka event: %w", err)
	}
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topicName,
		Value: sarama.StringEncoder(encoded),
	}

	return nil
}
