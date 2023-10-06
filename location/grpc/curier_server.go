package grpc

import (
	"context"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourierServer struct {
	CourierLocationRepository domain.CourierLocationRepositoryInterface
	pb.UnsafeCourierServer
}

func (courierServer CourierServer) GetCourierLatestPosition(ctx context.Context, req *pb.GetCourierLatestPositionRequest) (*pb.GetCourierLatestPositionResponse, error) {
	courierLatestPosition, err := courierServer.CourierLocationRepository.GetLatestPositionCourierById(ctx, req.CourierId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return &pb.GetCourierLatestPositionResponse{
		Latitude:  courierLatestPosition.Latitude,
		Longitude: courierLatestPosition.Longitude,
	}, err
}
