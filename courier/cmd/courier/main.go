package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	"github.com/Sharykhin/go-delivery-dymas/courier/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/courier/grpc"
	"github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/courier/postgres"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
	pborder "github.com/Sharykhin/go-delivery-dymas/proto/generate/order/v1"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg.Add(2)
	go runHttpServer(ctx, config, &wg, courierService)
	go runGRPC(ctx, config, &wg, courierRepository)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, courierService domain.CourierServiceInterface) {
	defer wg.Done()
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
	pkghttp.RunServer(ctx, router, ":"+config.PortServerCourier)
}

func runGRPC(ctx context.Context, config env.Config, wg *sync.WaitGroup, repo domain.CourierRepositoryInterface) {
	defer wg.Done()
	lis, err := net.Listen("tcp", config.CourierAssignerGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	courierAssignServer := grpc.NewServer()
	pborder.RegisterAssignCourierServer(courierAssignServer, couriergrpc.NewAssignCourierServer(repo))
	go func() {
		if err := courierAssignServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
	<-ctx.Done()
	courierAssignServer.GracefulStop()
}
