package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hash/fnv"
	"time"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderRepository needs for managing order.
type OrderRepository struct {
	client *sql.DB
}

// SaveOrder saves orders in db.
func (repo *OrderRepository) SaveOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	query := "insert into orders (customer_phone_number, created_at, status) values ($1, $2, $3) RETURNING id, customer_phone_number, status, created_at"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		order.CustomerPhoneNumber,
		order.CreatedAt,
		order.Status,
	)

	var orderRow domain.Order

	err := row.Scan(&orderRow.ID, &orderRow.CustomerPhoneNumber, &orderRow.Status, &orderRow.CreatedAt)

	return &orderRow, err
}

// ChangeOrderStatusAfterValidation Change order status  when taking validation.
func (repo *OrderRepository) ChangeOrderStatusAfterValidation(
	ctx context.Context,
	orderID string,
	orderValidation domain.OrderValidation,
) (orderRow *domain.Order, err error) {
	tx, err := repo.client.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func(tx *sql.Tx) (err error) {
		if err != nil {
			tx.Rollback()

			return
		}

		err = tx.Rollback()

		if errors.Is(err, sql.ErrTxDone) {
			err = nil
		}
		return
	}(tx)
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock($1)", repo.hashOrderID(orderID))
	if err != nil {
		return
	}

	query := "SELECT courier, courier_error FROM order_validations WHERE order_id=$1"

	row := tx.QueryRowContext(
		ctx,
		query,
		orderID,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	err = row.Scan(&orderValidation.Courier, &orderValidation.CourierError)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	query = "INSERT INTO order_validations (order_id, courier, courier_error, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (order_id)" +
		" DO UPDATE SET courier = $2 "
	_, err = tx.ExecContext(
		ctx,
		query,
		orderID,
		orderValidation.Courier,
		orderValidation.CourierError,
		time.Now(),
	)

	if !orderValidation.CheckValidation() && err == nil {
		err = tx.Commit()

		return
	}

	if err != nil {
		return
	}

	query = "Update orders SET status = $1 WHERE id = $2 RETURNING id, customer_phone_number, status, created_at"

	row = tx.QueryRowContext(
		ctx,
		query,
		domain.OrderStatusAccepted,
		orderID,
	)

	orderRow = &domain.Order{}

	err = row.Scan(&orderRow.ID, &orderRow.CustomerPhoneNumber, &orderRow.Status, &orderRow.CreatedAt)

	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (repo *OrderRepository) hashOrderID(orderID string) int64 {
	h := fnv.New64a()
	h.Write([]byte(orderID))
	return int64(h.Sum64())
}

// GetOrderByID gets order status and order id from database by uuid order and return model with order id.
func (repo *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	query := "SELECT id, status FROM orders WHERE id=$1 FOR SHARE"
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
