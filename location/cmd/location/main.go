package main

import (
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"log"
	"os"
	"strconv"
)

func main() {
	address := "localhost:9092"
	if os.Getenv("KAFKA_BROKERS") != "" {
		address = os.Getenv("KAFKA_BROKERS")
	}
	host := "localhost"
	port := "6379"
	db := 0
	var addr string
	addr = host + ":" + port
	if os.Getenv("REDIS_HOST") != "" && os.Getenv("REDIS_PORT") != "" {
		host = os.Getenv("REDIS_HOST")
		port = os.Getenv("REDIS_PORT")
		addr = host + ":" + port
	}

	if os.Getenv("DB_REDIS") != "" {
		dbConvert, error := strconv.Atoi(os.Getenv("DB_REDIS"))
		db = dbConvert
		if error != nil {
			panic(error)
		}
	}
	portServer := ":8081"

	if os.Getenv("HTTP_PORT") != "" {
		portServer = os.Getenv("HTTP_PORT")
	}
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
