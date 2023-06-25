package main

import (
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}

	publisher, err := kafka.NewPublisher(sarama.NewConfig(), config.Address)
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}
	redisClient := redis.CreateConnect(config.Addr, config.Db)
	repo := redis.CreateCouriersRepository(redisClient)
	courierService := kafka.NewCourierService(publisher, repo)
	locationHandler := handler.NewLocationHandler(courierService)
	router := http.NewRouter().CreateRouter(locationHandler, mux.NewRouter())
	if err := http.RunServer(router, ":"+config.PortServer); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
