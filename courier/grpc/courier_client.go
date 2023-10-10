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

func NewLocationClient(locationConnection *grpc.ClientConn, repo domain.CourierRepositoryInterface) *CourierLocationPositionClient {
	clientCourier := pb.NewCourierClient(locationConnection)
	return &CourierLocationPositionClient{
		courierClientGrpc: clientCourier,
	}
}
func (cl CourierLocationPositionClient) GetCourierLatestPosition(ctx context.Context, courierID string) (*domain.LocationPosition, error) {
	courierLatestPositionResponse, err := cl.courierClientGrpc.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
	locationPosition := domain.LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}

	if err != nil {
		return nil, err
	}

	return &locationPosition, nil
}
func NewCourierConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
