package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	kafakconsumer "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
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

	locationConsumer := kafka.NewCourierLocationConsumer(repo)
	err := kafakconsumer.NewConsumer[domain.CourierLocation]("latest_position_courier").RegisterJSONHandler(ctx, locationConsumer).StartConsuming(ctx)
}
