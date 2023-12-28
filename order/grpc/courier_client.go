package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
)

// CourierClient provides client for communicate with grpc server
type CourierClient struct {
	courierClientGRPC pb.CourierClient
}

// NewCourierClient creates new courier client for communicate with server by grpc
func NewCourierClient(courierConnection *grpc.ClientConn) *CourierClient {
	clientCourier := pb.NewCourierClient(courierConnection)

	return &CourierClient{
		courierClientGRPC: clientCourier,
	}
}

// GetFirstAvailableCourier gets first courier from courier service
func (cl *CourierClient) GetFirstAvailableCourier(ctx context.Context, courierID string) (*domain.Order, error) {
	courierLatestPositionResponse, err := cl.courierClientGRPC.GetFirstAvailableCourier(ctx, &pb.Null)
	code, ok := status.FromError(err)

	if ok && code.Code() == codes.NotFound {
		log.Printf("Not Found: %v\n", err)

		return nil, domain.ErrCourierNotFound
	}

	if err != nil {
		return nil, err
	}

	locationPosition := domain.Order{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}

	return &locationPosition, nil
}

// NewCourierConnection gets courier connection use grpc protocol
func NewCourierConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
