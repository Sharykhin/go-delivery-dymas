package main

import (
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/courier/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/courier/grpc"
	"github.com/Sharykhin/go-delivery-dymas/courier/http"
	"github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/courier/postgres"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	clientPostgres, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	defer clientPostgres.Close()
	courierConnection, err := couriergrpc.NewCourierConnection(config.CourierGrpcAddress)
	if err != nil {
		log.Panicf("Error Courier Server Connection: %v\n", err)
	}
	defer courierConnection.Close()
	courierClient := pb.NewCourierClient(courierConnection)
	courierRepository := postgres.NewCourierRepository(clientPostgres)
	courierHandler := handler.NewCourierHandler(courierRepository, courierClient)
	courierLatestPositionUrl := fmt.Sprintf("/couriers/{id:%s}", http.UuidRegexp)
	routes := map[string]http.Route{"/couriers": {
		Handler: courierHandler.HandlerCourierCreate,
		Method:  "POST",
	},
		courierLatestPositionUrl: {
			Handler: courierHandler.HandlerGetCourierLatestPosition,
			Method:  "GET",
		},
	}
	router := http.NewCourierRoute(routes, mux.NewRouter())
	if err := http.RunServer(router, ":"+config.PortServerCourier); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
