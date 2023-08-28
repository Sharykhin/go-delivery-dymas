package main

import (
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
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
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	clientPostgres, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	defer clientPostgres.Close()
	redisClient := redis.NewConnect(config.RedisAddress, config.Db)
	courierLocationRepository := redis.NewCourierLocationRepository(redisClient)
	courierLocationService := domain.NewCourierLocationService(courierLocationRepository, publisher)
	courierRepository := postgres.NewCourierRepository(clientPostgres)
	locationHandler := handler.NewLocationHandler(courierLocationService)
	courierCreateHandler := handler.NewCourierCreateHandler(courierRepository)
	router := http.NewRouteCourierLocation().NewCourierLocationRoute(locationHandler, mux.NewRouter())
	http.NewCourierCreateRoute().NewCourierCreateRoute(courierCreateHandler, router)
	if err := http.RunServer(router, ":"+config.PortServer); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
