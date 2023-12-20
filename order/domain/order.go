package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrOrderNotFound = errors.New("order was not found")

// Order is a model of an order.
type Order struct {
	ID                  string    `json:"id"`
	CourierID           string    `json:"courier_id"`
	CustomerPhoneNumber string    `json:"customer_phone_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// OrderPublisherInterface publish message some systems.
type OrderPublisherInterface interface {
	PublishOrder(ctx context.Context, order *Order) error
}

// OrderClientInterface save order.
type OrderClientInterface interface {
	SaveOrder(ctx context.Context, phoneNumber string) (*Order, error)
}

// OrderService provides information about order and save order
type OrderService struct {
	orderRepository OrderRepositoryInterface
	orderPublisher  OrderPublisherInterface
}

// OrderRepositoryInterface saves and reads courier from storage
type OrderRepositoryInterface interface {
	SaveOrder(ctx context.Context, order *Order) (*Order, error)
	GetStatusByOrderId(ctx context.Context, orderID string) (*Order, error)
}

// OrderServiceInterface gets information about courier and latest position courier from storage
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, courierID string) (*Order, error)
	GetStatusByOrderId(ctx context.Context, orderID string) (*Order, error)
}

// CreateOrder  new order in db
func (s *OrderService) CreateOrder(ctx context.Context, order *Order) (*Order, error) {

	order, err := s.orderRepository.SaveOrder(
		ctx,
		order,
	)

	s.orderPublisher.PublishOrder(ctx, order)

	if err != nil {
		return nil, fmt.Errorf("failed to create order in database: %w", err)
	}

	return order, nil
}

// GetStatusByOrderId returns status and order id data
func (s *OrderService) GetStatusByOrderId(ctx context.Context, orderId string) (*Order, error) {

	order, err := s.orderRepository.GetStatusByOrderId(
		ctx,
		orderId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create order in database: %w", err)
	}

	return order, nil
}

// NewOrderService creates new order service
func NewOrderService(repo OrderRepositoryInterface, orderPublisher OrderPublisherInterface) *OrderService {
	return &OrderService{
		orderRepository: repo,
		orderPublisher:  orderPublisher,
	}
}

// NewOrder creates new order for saving in db
func NewOrder(phoneNumber string) *Order {
	return &Order{
		CustomerPhoneNumber: phoneNumber,
		CreatedAt:           time.Now(),
		Status:              "pending",
	}
}
