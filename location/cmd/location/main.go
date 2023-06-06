package main

import (
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
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
	address := env.GetKafkaConfig()
	var addr string
	addr, db := env.GetRedisConfig()
	portServer := env.GetServerEnv()
	publisher := kafka.NewPublisher(sarama.NewConfig(), address)
	client := redis.CreateConnect(addr, db)
	repo := redis.CreateCouriersRepository(client)
	courierService := domain.NewCourierService(publisher, repo)
	locationHandler := handler.NewLocationHandler(courierService)
	router := http.NewRouter().CreateRouter(locationHandler, mux.NewRouter())
	if err := http.RunServer(router, portServer); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
