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

// GetFirstAvailableCourier gets first courier available.
func (courierServer CourierServer) GetFirstAvailableCourier(ctx context.Context, req *pb.Null) (*pb.GetAvailableCourierResponse, error) {
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

	return &pb.GetAvailableCourierResponse{
		CourierId: courier.ID,
	}, err
}
