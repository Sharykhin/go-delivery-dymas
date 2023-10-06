package grpc

import (
	"context"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc"
)

type LocationPositionService struct {
	CourierClient pb.CourierClient
}

func (s LocationPositionService) GetCourierLatestPosition(ctx context.Context, courierID string) (*domain.LocationPosition, error) {
	courierLatestPositionResponse, err := s.CourierClient.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
	if err != nil {
		return nil, err
	}
	locationPosition := domain.LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}
	return &locationPosition, nil
}

func NewLocationPositionService(locationConnection *grpc.ClientConn) *LocationPositionService {
	clientCourier := pb.NewCourierClient(locationConnection)
	return &LocationPositionService{
		CourierClient: clientCourier,
	}
}
