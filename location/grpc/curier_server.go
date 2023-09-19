package grpc

import (
	"context"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourierServer struct {
	courierLocationRepository domain.CourierLocationRepositoryInterface
}

func (courierServer CourierServer)GetCourierLatestPosition(ctx context.Context, getCourierLatestPositionRequest *pb.GetCourierLatestPositionRequest) (*pb.GetCourierLatestPositionResponse, error){
	courierLatestPosition, err := courierServer.courierLocationRepository.GetLatestPositionCourierById(ctx, getCourierLatestPositionRequest.CourierId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return pb.GetCourierLatestPositionResponse{
		Latitude: courierLatestPosition.Latitude,
		Longitude: courierLatestPosition.Longitude,
	}, err
}

func RunCourierServer(courierLocationRepository domain.CourierLocationRepositoryInterface, courCourierGrpcAddress string)  {
	lis, err := net.Listen("tcp", courCourierGrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	courierServer := grpc.NewServer()
	pb.RegisterCourierServer(courierServer, &CourierServer{
		courierLocationRepository: courierLocationRepository,
	})
	if err := courierServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}