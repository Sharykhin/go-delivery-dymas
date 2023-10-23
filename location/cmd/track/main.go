package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
	"log"
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
	consumerGroup, err := kafka.NewCourierLocationConsumer(repo, config.KafkaAddress, config.Verbose, config.Oldest, config.Assignor, "latest_position_courier")
	if err != nil {
		log.Panicf("Failed to create kafka consumer group: %v\n", err)
	}
	err = consumerGroup.ConsumeCourierLatestCourierGeoPositionMessage(ctx)
	if err != nil {
		log.Panicf("Failed to consume message: %v\n", err)
	}
}
