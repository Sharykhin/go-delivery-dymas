package grpc

import (
	"context"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CourierLocationPositionClient struct {
	courierClientGRPC pb.CourierClient
}

func NewCourierClient(locationConnection *grpc.ClientConn) *CourierLocationPositionClient {
	clientCourier := pb.NewCourierClient(locationConnection)
	return &CourierLocationPositionClient{
		courierClientGRPC: clientCourier,
	}
}
func (cl CourierLocationPositionClient) GetLatestPosition(ctx context.Context, courierID string) (*domain.LocationPosition, error) {
	courierLatestPositionResponse, err := cl.courierClientGRPC.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
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
