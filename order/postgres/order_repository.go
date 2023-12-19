package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderRepository needs for managing order.
type OrderRepository struct {
	client *sql.DB
}

// SaveOrder saves orders in db.
func (repo *OrderRepository) SaveOrder(ctx context.Context, phoneNumber string) (*domain.Order, error) {
	query := "insert into orders (phoneNumber) values ($1) ON CONFLICT DO NOTHING RETURNING id, courier_id, customer_phone_number, status, created_at"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		phoneNumber,
	)

	var orderRow domain.Order

	err := row.Scan(&orderRow.ID, &orderRow.CourierID, &orderRow.CustomerPhoneNumber, &orderRow.Status, &orderRow.CreatedAt)

	return &orderRow, err
}

// GetStatusByOrderId get order status and order id from database by uuid order and return model with order id.
func (repo *OrderRepository) GetStatusByOrderId(ctx context.Context, orderID string) (*domain.Order, error) {
	query := "SELECT *  FROM orders WHERE id=$1 FOR SHARE"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		orderID,
	)

	var orderRow domain.Order
	err := row.Scan(&orderRow.ID, &orderRow.Status)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}

	return &orderRow, err
}

// NewOrderRepository creates new order repository.
func NewOrderRepository(client *sql.DB) *OrderRepository {
	orderRepository := OrderRepository{
		client: client,
	}

	return &orderRepository
}
