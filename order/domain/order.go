package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const orderNewStatus = "pending"

// ErrOrderNotFound shows type this error, when we don't have order in db
var ErrOrderNotFound = errors.New("order was not found")

// Order is a model of an order.
type Order struct {
	ID                  string    `json:"id"`
	CourierID           string    `json:"courier_id"`
	CustomerPhoneNumber string    `json:"customer_phone_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// OrderPublisher publish message some systems.
type OrderPublisher interface {
	PublishOrder(ctx context.Context, order *Order) error
}

// OrderService provides information about order and save order
type OrderService struct {
	orderRepository OrderRepository
	orderPublisher  OrderPublisher
}

// OrderRepository OrderRepositoryInterface saves and reads courier from storage
type OrderRepository interface {
	SaveOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
}

// OrderServiceInterface gets information about courier and latest position courier from storage
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, courierID string) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
}

// CreateOrder creates new order and saves it in repository, and then publishes the corresponding event.
func (s *OrderService) CreateOrder(ctx context.Context, order *Order) (*Order, error) {

	order, err := s.orderRepository.SaveOrder(
		ctx,
		order,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create order in database: %w", err)
	}

	s.orderPublisher.PublishOrder(ctx, order)

	if err != nil {
		return nil, fmt.Errorf("failed to publish order: %w", err)
	}

	return order, nil
}

// GetOrderByID returns status and order id data
func (s *OrderService) GetOrderByID(ctx context.Context, orderId string) (*Order, error) {

	order, err := s.orderRepository.GetOrderByID(
		ctx,
		orderId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order by id: %w", err)
	}

	return order, nil
}

// NewOrderService creates new order service
func NewOrderService(orderRepo OrderRepository, orderPublisher OrderPublisher) *OrderService {
	return &OrderService{
		orderRepository: orderRepo,
		orderPublisher:  orderPublisher,
	}
}

// NewOrder creates new order for saving in db
func NewOrder(phoneNumber string) *Order {
	return &Order{
		CustomerPhoneNumber: phoneNumber,
		CreatedAt:           time.Now(),
		Status:              orderNewStatus,
	}
}
