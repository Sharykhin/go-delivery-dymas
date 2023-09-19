package main

import (
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	"github.com/Sharykhin/go-delivery-dymas/location/grpc"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"log"
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
	grpc.RunCourierServer(repo, config.CourierGrpcAddress)
}