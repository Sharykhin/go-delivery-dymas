package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	"github.com/Sharykhin/go-delivery-dymas/courier/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/courier/grpc"
	"github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/courier/kafka"
	"github.com/Sharykhin/go-delivery-dymas/courier/postgres"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
)

func main() {
	var wg sync.WaitGroup
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}

	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	clientPostgres, err := sql.Open("postgres", connPostgres)

	if err != nil {
		log.Panicf("error connection database: %v\n", err)
	}

	defer clientPostgres.Close()

	courierLocationGRPCConnection, err := couriergrpc.NewCourierLocationPositionConnection(config.CourierGrpcAddress)

	if err != nil {
		log.Panicf("error courier gRPC client connection: %v\n", err)
	}
	defer courierLocationGRPCConnection.Close()

	courierRepository := postgres.NewCourierRepository(clientPostgres)
	publisher, err := pkgkafka.NewPublisher([]string{config.KafkaAddress}, []string{config.KafkaSchemaRegistryAddress}, kafka.OrderTopicValidation)
	if err != nil {
		log.Panicf("failed to create publisher: %v\n", err)
	}
	orderValidationPublisher := kafka.NewOrderValidationPublisher(publisher)
	courierLocationClient := couriergrpc.NewCourierLocationPositionClient(courierLocationGRPCConnection)
	courierServiceManager := domain.NewCourierServiceManager(courierRepository, orderValidationPublisher, courierLocationClient)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg.Add(2)
	go runHttpServer(ctx, config, &wg, courierServiceManager)
	go runOrderConsumer(ctx, courierServiceManager, &wg, config)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, courierService domain.CourierService) {
	defer wg.Done()
	courierHandler := handler.NewCourierHandler(courierService, pkghttp.NewHandler())
	courierLatestPositionURL := fmt.Sprintf(
		"/couriers/{id:%s}",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	routes := map[string]pkghttp.Route{"/couriers": {
		Handler: courierHandler.HandlerCourierCreate,
		Method:  []string{"POST"},
	},
		courierLatestPositionURL: {
			Handler: courierHandler.GetCourier,
			Method:  []string{"GET"},
		},
	}
	router := pkghttp.NewRoute(routes, mux.NewRouter())
	pkghttp.RunServer(ctx, router, ":"+config.PortServerCourier)
}

func runOrderConsumer(ctx context.Context, courierService domain.CourierService, wg *sync.WaitGroup, config env.Config) {
	defer wg.Done()
	orderConsumer := kafka.NewOrderConsumer(courierService)
	consumer, err := pkgkafka.NewConsumer(
		orderConsumer,
		config.KafkaAddress,
		config.Verbose,
		config.Oldest,
		config.Assignor,
		kafka.OrderTopic,
		[]string{config.KafkaSchemaRegistryAddress},
	)

	if err != nil {
		log.Panicf("Failed to create kafka consumer group: %v\n", err)
	}

	err = consumer.ConsumeMessage(ctx)

	if err != nil {
		log.Panicf("Failed to consume message: %v\n", err)
	}
}
