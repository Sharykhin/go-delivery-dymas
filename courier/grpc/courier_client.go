package grpc

import (
	"context"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CourierLocationPositionClient struct {
	courierClientGrpc pb.CourierClient
}

func NewNewCourierClient(locationConnection *grpc.ClientConn) *CourierLocationPositionClient {
	clientCourier := pb.NewCourierClient(locationConnection)
	return &CourierLocationPositionClient{
		courierClientGrpc: clientCourier,
	}
}
func (cl CourierLocationPositionClient) GetCourierLatestPosition(ctx context.Context, courierID string) (*domain.LocationPosition, error) {
	courierLatestPositionResponse, err := cl.courierClientGrpc.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
	if err != nil {
		return nil, err
	}
	locationPosition := domain.LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}

	return &locationPosition, nil
}
func NewCourierConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
