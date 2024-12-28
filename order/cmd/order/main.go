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

	"github.com/Sharykhin/go-delivery-dymas/order/kafka"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	"github.com/Sharykhin/go-delivery-dymas/order/env"
	"github.com/Sharykhin/go-delivery-dymas/order/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/order/postgres"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
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
		log.Panicf("Error connection database: %v\n", err)
	}

	defer clientPostgres.Close()
	orderRepo := postgres.NewOrderRepository(clientPostgres)
	publisher, err := pkgkafka.NewPublisher([]string{config.KafkaAddress}, []string{config.KafkaSchemaRegistryAddress}, kafka.OrderTopic)
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}

	orderPublisher := kafka.NewOrderPublisher(publisher)
	orderServiceManager := domain.NewOrderServiceManager(orderRepo, orderPublisher)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	wg.Add(2)
	go runHttpServer(ctx, config, &wg, orderServiceManager)
	go runOrderConsumer(ctx, orderServiceManager, &wg, config)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, orderService domain.OrderService) {
	defer wg.Done()

	orderHandler := handler.NewOrderHandler(orderService, pkghttp.NewHandler())
	orderURL := "/orders/{order_id}"
	orderCancelUrl := "/orders/cancel/{order_id}"

	routes := map[string]pkghttp.Route{
		"/orders": {
			Handler: orderHandler.HandleOrderCreate,
			Methods: []string{"POST"},
		},
		orderURL: {
			Handler: orderHandler.HandleGetByOrderID,
			Methods: []string{"GET"},
		},
		orderCancelUrl: {
			Handler: orderHandler.HandleOrderCancel,
			Methods: []string{"POST"},
		},
	}

	router := pkghttp.NewRoute(routes, mux.NewRouter())
	pkghttp.RunServer(ctx, router, ":"+config.PortServerOrder)
}

func runOrderConsumer(ctx context.Context, orderService domain.OrderService, wg *sync.WaitGroup, config env.Config) {
	defer wg.Done()
	orderConsumer := kafka.NewOrderConsumerValidation(orderService)
	consumer, err := pkgkafka.NewConsumer(
		orderConsumer,
		config.KafkaAddress,
		config.Verbose,
		config.Oldest,
		config.Assignor,
		kafka.OrderValidationsTopic,
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
