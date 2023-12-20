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

	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBName)
	clientPostgres, err := sql.Open("postgres", connPostgres)

	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}

	defer clientPostgres.Close()
	repoOrderPostgres := postgres.NewOrderRepository(clientPostgres)

	courierService := domain.NewOrderService(repoOrderPostgres)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	wg.Add(1)
	go runHttpServer(ctx, config, &wg, courierService)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, orderService *domain.OrderService) {

	orderHandler := handler.NewOrderHandler(orderService, pkghttp.NewHandler())
	orderURL := fmt.Sprintf(
		"/orders/{order_id:%s}",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)
	routes := map[string]pkghttp.Route{
		"/orders": {
			Handler: orderHandler.HandlerOrderCreate,
			Method:  "POST",
		},
		orderURL: {
			Handler: orderHandler.HandlerOrderGetStatusByOrderId,
			Method:  "GET",
		},
	}

	router := pkghttp.NewRoute(routes, mux.NewRouter())
	pkghttp.RunServer(ctx, router, ":"+config.PortServerOrder)
	wg.Done()
}
