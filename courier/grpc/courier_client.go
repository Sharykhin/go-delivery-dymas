package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
)

// CourierLocationPositionClient provides client for communicate with grpc server
type CourierLocationPositionClient struct {
	courierLocationClientGRPC pb.CourierLocationClient
}

// NewCourierLocationClient creates new courier client for communicate with server by grpc
func NewCourierLocationClient(locationConnection *grpc.ClientConn) *CourierLocationPositionClient {
	courierLocationClient := pb.NewCourierLocationClient(locationConnection)

	return &CourierLocationPositionClient{
		courierLocationClientGRPC: courierLocationClient,
	}
}

// GetLatestPosition gets latest position from courier server
func (cl *CourierLocationPositionClient) GetLatestPosition(ctx context.Context, courierID string) (*domain.LocationPosition, error) {
	courierLatestPositionResponse, err := cl.courierLocationClientGRPC.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
	code, ok := status.FromError(err)

	if ok && code.Code() == codes.NotFound {
		log.Printf("Not Found: %v\n", err)
		return nil, domain.ErrCourierNotFound
	}

	if err != nil {
		return nil, err
	}

	locationPosition := domain.LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}

	return &locationPosition, nil
}

// NewCourierLocationConnection gets courier connection use grpc protocol
func NewCourierLocationConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
