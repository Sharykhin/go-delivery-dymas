package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrCourierNotFound shows type this error, when we don't have courier in db
var ErrCourierNotFound = errors.New("courier was not found")

// OrderValidationPublisher publish order validation message in queue for order service.
type OrderValidationPublisher interface {
	PublishValidationResult(ctx context.Context, courierAssignment *CourierAssignment) error
}

// CourierWithLatestPosition is a model of a courier, which provides general information and the latest courier position.
type CourierWithLatestPosition struct {
	LatestPosition *LocationPosition
	FirstName      string
	ID             string
	IsAvailable    bool
}

// CourierLocationClient gets latest location position courier.
type CourierLocationClient interface {
	GetLatestPosition(ctx context.Context, courierID string) (*LocationPosition, error)
}

// LocationPosition describes a geo position with simple latitude and longitude coordinates. In the courier domain is it used in order to store the latest courier position.
type LocationPosition struct {
	Latitude  float64
	Longitude float64
}

// CourierServiceManager provides information about courier and latest position from storage
type CourierServiceManager struct {
	courierLocationClient    CourierLocationClient
	courierRepository        CourierRepositoryInterface
	orderValidationPublisher OrderValidationPublisher
}

// CourierRepositoryInterface saves and reads courier from storage
type CourierRepositoryInterface interface {
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
	AssignOrderToCourier(ctx context.Context, orderID string) (CourierAssignment *CourierAssignment, err error)
	UnassignOrder(ctx context.Context, orderID string) error
}

// CourierAssignment has order assign courier
type CourierAssignment struct {
	OrderID   string    `json:"order_id"`
	CourierID string    `json:"courier_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CourierService gets information about courier and latest position courier from storage
type CourierService interface {
	GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierWithLatestPosition, error)
	AssignOrderToCourier(ctx context.Context, orderID string) error
	UnassignOrder(ctx context.Context, orderID string) error
	SaveCourier(ctx context.Context, courier *Courier) (*Courier, error)
}

// Courier provides information about courier
type Courier struct {
	ID          string
	FirstName   string
	IsAvailable bool
}

// AssignOrderToCourier assign order to courier and send message in queue
func (s *CourierServiceManager) AssignOrderToCourier(ctx context.Context, orderID string) error {

	courierAssigment, err := s.courierRepository.AssignOrderToCourier(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to save a courier assigments in the repository: %w", err)
	}

	err = s.orderValidationPublisher.PublishValidationResult(ctx, courierAssigment)

	if err != nil {
		return fmt.Errorf("failed to publish a order message validation in kafka: %w", err)
	}

	return nil
}

// UnassignOrder remove assigned order to courier and do available courier
func (s *CourierServiceManager) UnassignOrder(ctx context.Context, orderID string) error {

	err := s.courierRepository.UnassignOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to unassigment courier to order and do courier available: %w", err)
	}

	return nil
}

// GetCourierWithLatestPosition gets latest position from server and storage
func (s *CourierServiceManager) GetCourierWithLatestPosition(ctx context.Context, courierID string) (*CourierWithLatestPosition, error) {

	courier, err := s.courierRepository.GetCourierByID(
		ctx,
		courierID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get courier from the repository: %w", err)
	}

	courierLatestPositionResponse, err := s.courierLocationClient.GetLatestPosition(ctx, courierID)
	isErrCourierNotFound := errors.Is(err, ErrCourierNotFound)

	if err != nil && !isErrCourierNotFound {
		return nil, fmt.Errorf("failed to get courier latest position: %w", err)
	}

	var locationPosition *LocationPosition

	if courierLatestPositionResponse != nil {
		locationPosition = &LocationPosition{
			Latitude:  courierLatestPositionResponse.Latitude,
			Longitude: courierLatestPositionResponse.Longitude,
		}
	}

	courierResponse := CourierWithLatestPosition{
		FirstName:      courier.FirstName,
		ID:             courier.ID,
		IsAvailable:    courier.IsAvailable,
		LatestPosition: locationPosition,
	}

	return &courierResponse, nil
}

func (s *CourierServiceManager) SaveCourier(ctx context.Context, courier *Courier) (*Courier, error) {

	return s.courierRepository.SaveCourier(ctx, courier)
}

// NewCourierServiceManager creates new courier service manager
func NewCourierServiceManager(
	courierRepository CourierRepositoryInterface,
	orderValidationPublisher OrderValidationPublisher,
	courierLocationClient CourierLocationClient,
) *CourierServiceManager {
	return &CourierServiceManager{
		courierRepository:        courierRepository,
		orderValidationPublisher: orderValidationPublisher,
		courierLocationClient:    courierLocationClient,
	}
}
