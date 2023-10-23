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
	topicName string
	producer  sarama.AsyncProducer
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

// NewPublisher creates a new instance of Kafka publisher and returns general Publisher interface.
// As for given parameters it only accepts topic name. Managing producer Partitioner and RequiredAcks
// are not supported yet but will be implemented in latter versions.
func NewPublisher(topicName string) Publisher {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	address := os.Getenv("KAFKA_ADDRESS")
	if address == "" {
		log.Panic("env variable KAFKA_ADDRESS is required")
	}
	producer, err := sarama.NewAsyncProducer([]string{address}, config)
	if err != nil {
		log.Panicf("failed to initialize async producer: %v", err)
	}

	return &publisher{
		producer: producer,
	}
}
