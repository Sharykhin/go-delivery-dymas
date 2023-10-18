package main

import (
	"context"
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
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	wg.Add(2)
	go runHttpServer(config, wg)
	go runGRPC(config, wg)
	wg.Wait()
}

func runHttpServer(config env.Config, wg sync.WaitGroup) {

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
	http.RunServer(router, ":"+config.PortServer)
	redisClient.Close()
	wg.Done()
}

func runGRPC(config env.Config, wg sync.WaitGroup) {
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	client, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	repo := postgres.NewCourierLocationRepository(client)
	lis, err := net.Listen("tcp", config.CourierGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	courierLocationServer := grpc.NewServer()
	pb.RegisterCourierServer(courierLocationServer, &couriergrpc.CourierServer{
		CourierLocationRepository: repo,
	})
	go func() {
		if err := courierLocationServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	courierLocationServer.GracefulStop()
	client.Close()
	wg.Done()
}
