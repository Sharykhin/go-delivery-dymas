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

// GetCourierLatestPosition gets courier latest position.
func (courierServer CourierServer) GetCourierLatestPosition(ctx context.Context, req *pb.GetCourierLatestPositionRequest) (*pb.GetCourierLatestPositionResponse, error) {
	courierLatestPosition, err := courierServer.CourierLocationRepository.GetLatestPositionCourierByID(ctx, req.CourierId)

	isErrCourierNotFound := err != nil && errors.Is(err, domain.ErrCourierNotFound)
	if isErrCourierNotFound {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Position Not found: %v", err),
		)
	}

	if err != nil {
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
