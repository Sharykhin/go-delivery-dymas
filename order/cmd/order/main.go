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
	"github.com/Sharykhin/go-delivery-dymas/order/grpc"
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

	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBName)
	clientPostgres, err := sql.Open("postgres", connPostgres)

	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}

	defer clientPostgres.Close()
	orderRepo := postgres.NewOrderRepository(clientPostgres)
	publisher, err := pkgkafka.NewPublisher(config.KafkaAddress, "orders")
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}
	orderPublisher := kafka.NewOrderPublisher(publisher)
	orderService := domain.NewOrderService(orderRepo, orderPublisher)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	wg.Add(2)
	go runHttpServer(ctx, config, &wg, orderService)
	go runConsumer(orderRepo, orderPublisher, &wg, config)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, orderService *domain.OrderService) {
	defer wg.Done()

	orderHandler := handler.NewOrderHandler(orderService, pkghttp.NewHandler())
	orderURL := fmt.Sprintf(
		"/orders/{order_id:%s}",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	routes := map[string]pkghttp.Route{
		"/orders": {
			Handler: orderHandler.HandleOrderCreate,
			Method:  "POST",
		},
		orderURL: {
			Handler: orderHandler.HandleGetByOrderID,
			Method:  "GET",
		},
	}

	router := pkghttp.NewRoute(routes, mux.NewRouter())
	pkghttp.RunServer(ctx, router, ":"+config.PortServerOrder)
}

func runConsumer(orderRepository domain.OrderRepository, orderPublisher domain.OrderPublisher, wg *sync.WaitGroup, config env.Config) {
	defer wg.Done()
	courierConnection, err := grpc.NewCourierConnection(config.OrderGrpcAddress)
	if err != nil {
		log.Panicf("Failed to create kafka consumer group: %v\n", err)
	}
	courierAssignClient := grpc.NewAssignCourierClient(courierConnection)
	orderConsumer := kafka.NewOrderConsumer(orderRepository, orderPublisher, courierAssignClient)
	consumer, err := pkgkafka.NewConsumer(
		orderConsumer,
		config.KafkaAddress,
		config.Verbose,
		config.Oldest,
		config.Assignor,
		"orders",
	)

	if err != nil {
		log.Panicf("Failed to create kafka consumer group: %v\n", err)
	}

	ctx := context.Background()

	err = consumer.ConsumeMessage(ctx)

	if err != nil {
		log.Panicf("Failed to consume message: %v\n", err)
	}
}
