package grpc

import (
	"github.com/Sharykhin/go-delivery-dymas/courier/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
)

func main()  {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	conn, err := grpc.Dial(config.CourierGrpcAddress, opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewCourierClient(conn)

	return client
}
