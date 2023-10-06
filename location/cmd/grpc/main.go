package main

import (
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/location/grpc"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
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
}
