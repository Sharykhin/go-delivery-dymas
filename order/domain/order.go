package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const OrderNewStatus = "pending"
const EventOrderCreated = "created"
const EventOrderUpdated = "updated"
const OrderStatusAccepted = "accepted"
const OrderStatusCanceled = "canceled"

var ErrorTransactionNotBegin = errors.New("transaction not begin")

// ErrOrderNotFound shows type this error, when we don't have order in db
var ErrOrderNotFound = errors.New("order was not found")
var ErrOrderValidationNotFound = errors.New("order validation was not found")
var ErrorCanceledOrder = errors.New("order can not be canceled order has incorrect status for canceling")

// Order is a model of an order.
type Order struct {
	ID                  string    `json:"id"`
	CourierID           *string   `json:"courier_id"`
	CustomerPhoneNumber string    `json:"customer_phone_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// OrderValidation imagine entity for order validation for saving in db
type OrderValidation struct {
	OrderID            string
	CourierValidatedAt time.Time
	UpdatedAt          time.Time
	CourierError       string
}

// OrderValidationPayload imagine payload for order validation for different services
type OrderValidationPayload struct {
	CourierID string
}

// OrderPublisher publish message some systems.
type OrderPublisher interface {
	PublishOrder(ctx context.Context, order *Order, event string) error
}

// OrderServiceManager provides information about order and save order
type OrderServiceManager struct {
	orderRepository OrderRepository
	orderPublisher  OrderPublisher
}

// OrderRepository OrderRepositoryInterface saves and reads courier from storage
type OrderRepository interface {
	Commit(ctx context.Context) error
	BeginTx(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context)
	LockTransaction(ctx context.Context, orderID string) error
	SaveOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	SaveOrderValidation(ctx context.Context, orderValidation *OrderValidation) error
	UpdateOrder(ctx context.Context, order *Order) error
	GetOrderValidationByID(ctx context.Context, orderID string) (*OrderValidation, error)
	UpdateOrderValidation(ctx context.Context, orderValidation *OrderValidation) error
}

// OrderService gets information about courier and latest position courier from storage
type OrderService interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	CancelOrderByID(ctx context.Context, orderID string) error
	ValidateOrderForService(ctx context.Context, serviceName string, orderID string, orderValidationPayload *OrderValidationPayload) error
}

// CheckValidation checks validation for all services after that we change status order if order pass validation
func (orderValidation *OrderValidation) CheckValidation() bool {

	if orderValidation.CourierError != "" {
		return false
	}

	return true
}

// CreateOrder creates new order and saves it in repository, and then publishes the corresponding event.
func (s *OrderServiceManager) CreateOrder(ctx context.Context, order *Order) (*Order, error) {

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

// ValidateOrderForService updates order status and creates or saves order validation
func (s *OrderServiceManager) ValidateOrderForService(ctx context.Context, serviceName string, orderID string, orderValidationPayload *OrderValidationPayload) error {
	ctx, err := s.orderRepository.BeginTx(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer s.orderRepository.Rollback(ctx)

	err = s.orderRepository.LockTransaction(ctx, orderID)

	if err != nil {
		return fmt.Errorf("failed to lock transaction: %w", err)
	}

	defer s.orderRepository.Rollback(ctx)

	order, err := s.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	orderValidation, err := s.orderRepository.GetOrderValidationByID(ctx, orderID)

	createNewOrderValidation := errors.Is(err, ErrOrderValidationNotFound)

	if err != nil && !createNewOrderValidation {
		return fmt.Errorf("failed to get order validation: %w", err)
	}

	if orderValidation == nil {
		orderValidation = &OrderValidation{}
		orderValidation.OrderID = orderID
	}

	var isCourierUpdateInOrder bool

	switch serviceName {
	case "courier":
		order.CourierID = &orderValidationPayload.CourierID
		orderValidation.CourierValidatedAt = time.Now()
		isCourierUpdateInOrder = true
	}

	if createNewOrderValidation {
		err = s.orderRepository.SaveOrderValidation(
			ctx,
			orderValidation,
		)
	} else {
		err = s.orderRepository.UpdateOrderValidation(
			ctx,
			orderValidation,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to save order in database during validation: %w", err)
	}

	isOrderValidated := orderValidation.CheckValidation()
	if isOrderValidated {
		order.Status = OrderStatusAccepted
	}

	if isCourierUpdateInOrder || isOrderValidated {
		err = s.orderRepository.UpdateOrder(ctx, order)

		if err != nil {
			return fmt.Errorf("failed to order order in database during validation: %w", err)
		}

	}

	err = s.orderRepository.Commit(ctx)

	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	if isOrderValidated {
		err = s.orderPublisher.PublishOrder(ctx, order, EventOrderUpdated)

		if err != nil {
			return fmt.Errorf("failed to publish a order in the kafka: %w", err)
		}
	}

	return nil
}

// CancelOrderByID cancel order and publish in kafka message with event cancel for removing courier assign
func (s *OrderServiceManager) CancelOrderByID(ctx context.Context, orderID string) error {

	ctx, err := s.orderRepository.BeginTx(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer s.orderRepository.Rollback(ctx)

	err = s.orderRepository.LockTransaction(ctx, orderID)

	if err != nil {
		return fmt.Errorf("failed to lock transaction: %w", err)
	}

	order, err := s.orderRepository.GetOrderByID(
		ctx,
		orderID,
	)

	if err != nil {
		return fmt.Errorf("failed to retrieve order by id: %w", err)
	}

	if order.Status != OrderNewStatus {
		return ErrorCanceledOrder
	}

	order.Status = OrderStatusCanceled
	err = s.orderRepository.UpdateOrder(ctx, order)

	if err != nil {
		return fmt.Errorf("failed to cancel order by id: %w", err)
	}

	err = s.orderRepository.Commit(ctx)

	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	err = s.orderPublisher.PublishOrder(ctx, order, OrderStatusCanceled)

	if err != nil {
		return fmt.Errorf("failed to publish a order in the kafka: %w", err)
	}

	return nil
}

// GetOrderByID returns status and order id data
func (s *OrderServiceManager) GetOrderByID(ctx context.Context, orderID string) (*Order, error) {

	order, err := s.orderRepository.GetOrderByID(
		ctx,
		orderID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order by id: %w", err)
	}

	return order, nil
}

// NewOrderServiceManager creates new order service
func NewOrderServiceManager(
	orderRepo OrderRepository,
	orderPublisher OrderPublisher,
) *OrderServiceManager {
	return &OrderServiceManager{
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
