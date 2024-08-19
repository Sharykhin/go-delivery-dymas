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
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/location/grpc"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	wp "github.com/Sharykhin/go-delivery-dymas/location/workerpool"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
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
	repoPostgres := postgres.NewCourierLocationRepository(clientPostgres)
	publisher, err := pkgkafka.NewPublisher([]string{config.KafkaAddress}, []string{config.KafkaSchemaRegistryAddress}, kafka.LatestPositionCourierTopic)

	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}

	courierLocationPublisher := kafka.NewCourierLocationPublisher(publisher)
	redisClient := redis.NewConnect(config.RedisAddress, config.Db)
	defer redisClient.Close()

	repoRedis := redis.NewCourierLocationRepository(redisClient)
	courierService := domain.NewCourierLocationService(repoRedis, courierLocationPublisher)

	locationWorkerPool := wp.NewLocationPool(
		courierService,
		config.CourierLocationWorkerPoolCount,
		config.CourierLocationQueueSizeTasks,
		time.Duration(config.CourierLocationWorkerTimeoutGracefulShutdown),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	wg.Add(3)
	go locationWorkerPool.Run(ctx, &wg)
	go runHttpServer(ctx, config, &wg, locationWorkerPool)
	go runGRPC(ctx, config, &wg, repoPostgres)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, locationWorkerPool domain.CourierLocationWorkerPool) {

	locationHandler := handler.NewLocationHandler(locationWorkerPool, pkghttp.NewHandler())
	var courierLocationURL = fmt.Sprintf(
		"/courier/{courier_id:%s}/location",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)
	routes := map[string]pkghttp.Route{courierLocationURL: {
		Handler: locationHandler.HandlerCouriersLocation,
		Methods: []string{"POST"},
	},
	}
	router := pkghttp.NewRoute(routes, mux.NewRouter())
	pkghttp.RunServer(ctx, router, ":"+config.PortServer)
	wg.Done()
}

func runGRPC(ctx context.Context, config env.Config, wg *sync.WaitGroup, repo domain.CourierRepositoryInterface) {
	lis, err := net.Listen("tcp", config.CourierGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	courierLocationServer := grpc.NewServer()
	pb.RegisterCourierLocationPositionServer(courierLocationServer, &couriergrpc.CourierServer{
		CourierLocationRepository: repo,
	})
	go func() {
		if err := courierLocationServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
	<-ctx.Done()
	courierLocationServer.GracefulStop()
	wg.Done()
}
