package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourierServer struct {
	CourierLocationRepository domain.CourierRepositoryInterface
	pb.UnsafeCourierServer
}

// GetAssignCourier gets first courier available.
func (courierServer CourierServer) GetAssignCourier(ctx context.Context, req *pb.Empty) (*pb.GetFirstAvailableCourierResponse, error) {
	courier, err := courierServer.CourierLocationRepository.GetAppliedCourier(ctx)

	isErrCourierNotFound := err != nil && errors.Is(err, domain.ErrCourierNotFound)
	if isErrCourierNotFound {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("courier not found: %v", err),
		)
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("courier not found: %v", err),
		)
	}

	return &pb.GetAssignCourierResponse{
		CourierId: courier.ID,
	}, err
}
