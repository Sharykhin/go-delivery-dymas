package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/credentials/insecure"
)

func NewConnection(courierGrpcAddress string) pb.CourierClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(courierGrpcAddress, opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewCourierClient(conn)

	return client
}
