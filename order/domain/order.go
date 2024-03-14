package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const OrderNewStatus = "pending"
const EventOrderCreated = "created"
const EventOrderUpdated = "updated"
const OrderStatusAccepted = "accepted"

// ErrOrderNotFound shows type this error, when we don't have order in db
var ErrOrderNotFound = errors.New("order was not found")
var ErrOrderValidationNotFound = errors.New("order validation was not found")

// Order is a model of an order.
type Order struct {
	ID                  string    `json:"id"`
	CourierID           string    `json:"courier_id"`
	CustomerPhoneNumber string    `json:"customer_phone_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

type CourierPayload struct {
	CourierID string    `json:"courier_id"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderValidation imagine entity for order validation for saving in db
type OrderValidation struct {
	orderID           string
	CourierValidateAt time.Time
	UpdatedAt         time.Time
	CourierError      string
	serviceName       string
}

// OrderPublisher publish message some systems.
type OrderPublisher interface {
	PublishOrder(ctx context.Context, order *Order, event string) error
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
	SaveOrderValidation(
		ctx context.Context,
		orderID string,
		orderValidation *OrderValidation,
	) error
	UpdateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderValidationValidationById(ctx context.Context, orderID string) (*OrderValidation, error)
	UpdateOrderValidation(
		ctx context.Context,
		orderID string,
		orderValidation *OrderValidation,
	) (err error)
}

// OrderServiceInterface gets information about courier and latest position courier from storage
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	ChangeOrderStatusAfterValidation(ctx context.Context, serviceName string, orderID string, validationInfo []byte) error
}

func (orderValidation *OrderValidation) CheckValidation() bool {

	if orderValidation.CourierError != "" {
		return false
	}

	return true
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

	err = s.orderPublisher.PublishOrder(ctx, order, EventOrderCreated)

	if err != nil {
		return nil, fmt.Errorf("failed to publish order: %w", err)
	}

	return order, nil
}

// ChangeOrderStatusAfterValidation updates order status and creates or saves order validation
func (s *OrderService) ChangeOrderStatusAfterValidation(ctx context.Context, serviceName string, orderID string, validationInfo []byte) error {
	var orderValidation *OrderValidation
	var order = &Order{}

	var courierPayload CourierPayload

	orderValidation, err := s.orderRepository.GetOrderValidationValidationById(ctx, orderID)

	switch serviceName {
	case "courier":

		if err != nil && !errors.Is(err, ErrOrderValidationNotFound) {
			return fmt.Errorf("failed to get order validation: %v\n", err)
		}

		if err := json.Unmarshal(validationInfo, &courierPayload); err != nil {
			return fmt.Errorf("failed to unmarshal Kafka message into order struct: %v\n", err)
		}

		order.ID = orderID
		order.CourierID = courierPayload.CourierID
	}

	if errors.Is(err, ErrOrderValidationNotFound) {
		orderValidation = &OrderValidation{}
		err = s.orderRepository.SaveOrderValidation(
			ctx,
			orderID,
			orderValidation,
		)
	} else {
		err = s.orderRepository.UpdateOrderValidation(
			ctx,
			orderID,
			orderValidation,
		)
	}

	if err != nil && !errors.Is(err, ErrOrderValidationNotFound) {
		return err
	}

	order.ID = orderID
	if orderValidation.CheckValidation() {
		order.Status = OrderStatusAccepted
		order, err = s.orderRepository.UpdateOrder(ctx, order)

		if err != nil {
			return err
		}
	} else {
		order, err = s.orderRepository.UpdateOrder(ctx, order)

		if err != nil {
			return fmt.Errorf("failed to save a order in the repository: %w", err)
		}
	}

	err = s.orderPublisher.PublishOrder(ctx, order, EventOrderUpdated)

	if err != nil {
		return fmt.Errorf("failed to publish a order in the kafka: %w", err)
	}

	return nil
}

// GetOrderByID returns status and order id data
func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*Order, error) {

	order, err := s.orderRepository.GetOrderByID(
		ctx,
		orderID,
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
		Status:              OrderNewStatus,
	}
}
