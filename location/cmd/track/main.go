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
	repo, err := postgres.NewCourierLocationRepository(config.DbName, config.DbUser, config.PasswordDb)
	defer repo.Client.Close()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	consumerGroup, err := kafka.NewCourierLocationConsumer(repo, config.Address, config.Verbose, config.Oldest, config.Assignor)
	if err != nil {
		log.Println(err)
		return
	}
	err = consumerGroup.ConsumeCourierLatestCourierGeoPositionMessage(ctx)
	if err != nil {
		log.Println(err)
	}
}
