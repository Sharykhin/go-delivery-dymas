package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/env"
	couriergrpc "github.com/Sharykhin/go-delivery-dymas/location/grpc"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/kafka"
	"github.com/Sharykhin/go-delivery-dymas/location/postgres"
	"github.com/Sharykhin/go-delivery-dymas/location/redis"
	pkgkafka "github.com/Sharykhin/go-delivery-dymas/pkg/kafka"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	config, err := env.GetConfig()
	ctx := context.Background()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	clientPostgres, err := sql.Open("postgres", connPostgres)
	if err != nil {
		log.Panicf("Error connection database: %v\n", err)
	}
	//
	defer clientPostgres.Close()
	repoPostgres := postgres.NewCourierLocationRepository(clientPostgres)
	publisher, err := pkgkafka.NewPublisher(config.KafkaAddress, "latest_position_courier")
	if err != nil {
		log.Printf("failed to create publisher: %v\n", err)
		return
	}
	courierLocationPublisher := kafka.NewCourierLocationPublisher(publisher)
	redisClient := redis.NewConnect(config.RedisAddress, config.Db)
	defer redisClient.Close()
	repoRedis := redis.NewCourierLocationRepository(redisClient)
	courierService := domain.NewCourierLocationService(repoRedis, courierLocationPublisher)
	wg.Add(2)
	go runHttpServer(ctx, config, &wg, courierService)
	go runGRPC(ctx, config, &wg, repoPostgres)
	wg.Wait()
}

func runHttpServer(ctx context.Context, config env.Config, wg *sync.WaitGroup, courierService domain.CourierLocationServiceInterface) {

	locationHandler := handler.NewLocationHandler(courierService)
	router := http.NewRouter().NewRouter(locationHandler, mux.NewRouter())
	http.RunServer(ctx, router, ":"+config.PortServer)
	wg.Done()
}

func runGRPC(ctx context.Context, config env.Config, wg *sync.WaitGroup, repo domain.CourierRepositoryInterface) {
	lis, err := net.Listen("tcp", config.CourierGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	courierLocationServer := grpc.NewServer()
	pb.RegisterCourierServer(courierLocationServer, &couriergrpc.CourierServer{
		CourierLocationRepository: repo,
	})
	go func() {
		if err := courierLocationServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	courierLocationServer.GracefulStop()
	wg.Done()
}
