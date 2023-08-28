package main

import (
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}

	publisher, err := kafka.NewCourierLocationPublisher(config.KafkaAddress)
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}

	redisClient := redis.NewConnect(config.RedisAddress, config.Db)
	courierLocationRepository := redis.NewCourierLocationRepository(redisClient)
	courierLocationService := domain.NewCourierLocationService(courierLocationRepository, publisher)
	locationHandler := handler.NewLocationHandler(courierLocationService)
	router := http.NewRouteCourierLocation().NewCourierLocationRoute(locationHandler, mux.NewRouter())
	if err := http.RunServer(router, ":"+config.PortServer); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
