package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	"github.com/Sharykhin/go-delivery-dymas/courier/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/courier/grpc"
	"github.com/Sharykhin/go-delivery-dymas/courier/http"
	"github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/courier/postgres"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}

	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBName)
	clientPostgres, err := sql.Open("postgres", connPostgres)

	if err != nil {
		log.Panicf("error connection database: %v\n", err)
	}

	defer clientPostgres.Close()

	courierGRPCConnection, err := couriergrpc.NewCourierConnection(config.CourierGrpcAddress)

	if err != nil {
		log.Panicf("error courier gRPC client connection: %v\n", err)
	}
	defer courierGRPCConnection.Close()

	courierRepository := postgres.NewCourierRepository(clientPostgres)
	courierClient := couriergrpc.NewCourierClient(courierGRPCConnection)
	courierService := domain.NewCourierService(courierClient, courierRepository)
	courierHandler := handler.NewCourierHandler(courierService, pkghttp.NewHandler())
	courierLatestPositionURL := fmt.Sprintf(
		"/couriers/{id:%s}",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	routes := map[string]pkghttp.Route{"/couriers": {
		Handler: courierHandler.HandlerCourierCreate,
		Method:  "POST",
	},
		courierLatestPositionURL: {
			Handler: courierHandler.GetCourier,
			Method:  "GET",
		},
	}

	router := pkghttp.NewRoute(routes, mux.NewRouter())

	if err := http.RunServer(router, ":"+config.PortServerCourier); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
