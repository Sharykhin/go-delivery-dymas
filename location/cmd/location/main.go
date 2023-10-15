package main

import (
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/location/grpc"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	serverSignalErrChan := make(chan int)
	go runHttpServer(config, serverSignalErrChan)
	go runGRPC(config, serverSignalErrChan)
	<-serverSignalErrChan
}

func runHttpServer(config env.Config, serverSignalErrChan chan int) {

	publisher, err := kafka.NewCourierLocationPublisher(config.KafkaAddress)
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}
	redisClient := redis.NewConnect(config.RedisAddress, config.Db)
	repo := redis.NewCourierLocationRepository(redisClient)
	courierService := domain.NewCourierLocationService(repo, publisher)
	locationHandler := handler.NewLocationHandler(courierService)
	router := http.NewRouter().NewRouter(locationHandler, mux.NewRouter())
	if err := http.RunServer(router, ":"+config.PortServer); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
	serverSignalErrChan <- 1
}

func runGRPC(config env.Config, serverSignalErrChan chan int) {
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	client, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	defer client.Close()
	repo := postgres.NewCourierLocationRepository(client)
	lis, err := net.Listen("tcp", config.CourierGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	courierLocationServer := grpc.NewServer()
	pb.RegisterCourierServer(courierLocationServer, &couriergrpc.CourierServer{
		CourierLocationRepository: repo,
	})
	if err := courierLocationServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	serverSignalErrChan <- 1
}
