package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"log"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
	}
	ctx := context.Background()
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.PasswordDb, config.DbName)
	client, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
	}
	repo, err := postgres.NewCourierLocationRepository(client)
	defer repo.Client.Close()
	if err != nil {
		log.Panicf("failed to parse variable env: %v\n", err)
	}
	consumerGroup, err := kafka.NewCourierLocationConsumer(repo, config.Address, config.Verbose, config.Oldest, config.Assignor)
	if err != nil {
		log.Panic(err)
	}
	err = consumerGroup.ConsumeCourierLatestCourierGeoPositionMessage(ctx)
	if err != nil {
		log.Panic(err)
	}
}
