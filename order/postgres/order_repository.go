package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	courierPayload *domain.CourierPayload,
	statusValidation bool,
	orderStatusValidation string,
	serviceValidation string,
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

	query := "Update orders SET status = $1 WHERE id = $2 RETURNING id, customer_phone_number, status, created_at"

	row := tx.QueryRowContext(
		ctx,
		query,
		orderStatusValidation,
		courierPayload.OrderID,
	)

	orderRow = &domain.Order{}

	err = row.Scan(&orderRow.ID, &orderRow.CustomerPhoneNumber, &orderRow.Status, &orderRow.CreatedAt)

	if err != nil {
		return
	}

	query = fmt.Sprintf(
		"INSERT INTO order_validations (order_id, %s, created_at) VALUES ($1, $2, $3) ON CONFLICT (order_id) DO UPDATE SET %s = $2",
		serviceValidation,
		serviceValidation,
	)

	_, err = tx.ExecContext(
		ctx,
		query,
		courierPayload.OrderID,
		statusValidation,
		courierPayload.CreatedAt,
	)

	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
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
