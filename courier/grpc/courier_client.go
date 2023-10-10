package grpc

import (
	"context"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LocationPositionService struct {
	CourierClient     pb.CourierClient
	courierRepository domain.CourierRepositoryInterface
}

func (s LocationPositionService) GetCourierLatestPosition(ctx context.Context, courierID string) (*domain.CourierResponse, error) {
	courier, err := s.courierRepository.GetCourierByID(
		ctx,
		courierID,
	)
	if err != nil {
		return nil, err
	}
	courierLatestPositionResponse, err := s.CourierClient.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{CourierId: courierID})
	if err != nil {
		return nil, err
	}
	locationPosition := domain.LocationPosition{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}
	courierResponse := domain.CourierResponse{
		FirstName:      courier.FirstName,
		Id:             courier.Id,
		IsAvailable:    courier.IsAvailable,
		LatestPosition: &locationPosition,
	}
	return &courierResponse, nil
}

func NewLocationPositionService(locationConnection *grpc.ClientConn, repo domain.CourierRepositoryInterface) *LocationPositionService {
	clientCourier := pb.NewCourierClient(locationConnection)
	return &LocationPositionService{
		CourierClient:     clientCourier,
		courierRepository: repo,
	}
}

func NewCourierConnection(courierGrpcAddress string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return grpc.Dial(courierGrpcAddress, opts...)
}
