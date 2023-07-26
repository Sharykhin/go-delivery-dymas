package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"log"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	ctx := context.Background()
	repo, err := postgres.NewCourierLocationRepository()
	defer repo.Client.Close()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	consumerGroup, err := kafka.NewCourierLocationConsumer(repo, config.Address)
	if err != nil {
		log.Println(err)
		return
	}
	consumerGroup.ConsumeCourierLatestCourierGeoPositionMessage(ctx)
}
