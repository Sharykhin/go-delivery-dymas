package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/order/v1"
)

// AssignCourierClient AssignCourier provides client for communicate with grpc server
type AssignCourierClient struct {
	assignCourierGRPC pb.AssignCourierClient
}

// NewAssignCourierClient creates new assign courier client for communicate with server by grpc
func NewAssignCourierClient(assignCourierConnection *grpc.ClientConn) *AssignCourierClient {
	clientCourier := pb.NewAssignCourierClient(assignCourierConnection)

	return &AssignCourierClient{
		assignCourierGRPC: clientCourier,
	}
}

// GetAssignCourier gets first courier from courier service
func (ac *AssignCourierClient) GetAssignCourier(ctx context.Context, order domain.Order) (*domain.Order, error) {
	courier, err := ac.assignCourierGRPC.GetAssignCourier(ctx, &pb.Empty{})
	code, ok := status.FromError(err)

	if ok && code.Code() == codes.NotFound {
		log.Printf("Not Found: %v\n", err)

		return nil, domain.ErrCourierNotFound
	}

	if err != nil {
		return nil, err
	}

	order.CourierID = courier.CourierId

	return &order, nil
}

// NewCourierConnection gets courier connection use grpc protocol
func NewCourierConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
