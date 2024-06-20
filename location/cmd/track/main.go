package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Panicf("failed to parse variable env: %v\n", err)
	}
	ctx := context.Background()
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	client, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	defer client.Close()
	repo := postgres.NewCourierLocationRepository(client)
	courierLocationConsumer := kafka.NewCourierLocationConsumer(repo)
	consumer, err := pkgkafka.NewConsumer(
		courierLocationConsumer,
		config.KafkaAddress,
		config.Verbose,
		config.Oldest,
		config.Assignor,
		"latest_position_courier.v1",
		[]string{config.KafkaSchemaRegistryAddress},
	)

	if err != nil {
		log.Panicf("Failed to create kafka consumer group: %v\n", err)
	}

	err = consumer.ConsumeMessage(ctx)

	if err != nil {
		log.Panicf("Failed to consume message: %v\n", err)
	}
}
