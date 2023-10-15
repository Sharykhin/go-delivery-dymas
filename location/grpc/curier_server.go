package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourierServer struct {
	CourierLocationRepository domain.CourierRepositoryInterface
	pb.UnsafeCourierServer
}

func (courierServer CourierServer) GetCourierLatestPosition(ctx context.Context, req *pb.GetCourierLatestPositionRequest) (*pb.GetCourierLatestPositionResponse, error) {
	courierLatestPosition, err := courierServer.CourierLocationRepository.GetLatestPositionCourierById(ctx, req.CourierId)

	isErrorCourierNotFound := err != nil && errors.Is(err, domain.ErrorCourierNotFound)
	if isErrorCourierNotFound {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Position Not found: %v", err),
		)
	}
	if err != nil && !isErrorCourierNotFound {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Position Not found: %v", err),
		)
	}

	return &pb.GetCourierLatestPositionResponse{
		Latitude:  courierLatestPosition.Latitude,
		Longitude: courierLatestPosition.Longitude,
	}, err
}
