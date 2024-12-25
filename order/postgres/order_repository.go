package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
)

// OrderRepository needs for managing order.
type OrderRepository struct {
	client *sql.DB
}

type Client interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

// SaveOrder saves orders in db.
func (repo *OrderRepository) SaveOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {

	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "insert into orders (customer_phone_number, created_at, status) values ($1, $2, $3) RETURNING id, customer_phone_number, status, created_at"
	row := client.QueryRowContext(
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

// BeginTx Begin Transaction and set in context for using in repository methods
func (repo *OrderRepository) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := repo.client.BeginTx(ctx, nil)

	if err == nil {
		ctx = context.WithValue(ctx, "transaction", tx)
	}

	return ctx, err
}

// Rollback transaction in repository
func (repo *OrderRepository) Rollback(ctx context.Context) (err error) {
	tx := ctx.Value("transaction").(*sql.Tx)
	err = tx.Rollback()

	if errors.Is(err, sql.ErrTxDone) {
		err = nil

		return
	}

	return
}

// Commit transaction
func (repo *OrderRepository) Commit(ctx context.Context) (err error) {
	tx := ctx.Value("transaction").(*sql.Tx)

	return tx.Commit()
}

// UpdateOrder update order in db after get data from services.
func (repo *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "UPDATE orders SET status=$1, courier_id=$2 WHERE id = $3 RETURNING id, customer_phone_number, status, created_at, courier_id"
	_, err := client.ExecContext(
		ctx,
		query,
		order.Status,
		order.CourierID,
		order.ID,
	)

	return err
}

// GetOrderValidationByID GetOrderValidationValidationById gets order validation by id from db
func (repo *OrderRepository) GetOrderValidationByID(ctx context.Context, orderID string) (*domain.OrderValidation, error) {
	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "SELECT order_id, courier_validated_at, courier_error, updated_at FROM order_validations WHERE order_id=$1"

	row := client.QueryRowContext(
		ctx,
		query,
		orderID,
	)

	var orderValidation domain.OrderValidation

	err := row.Scan(&orderValidation.OrderID, &orderValidation.CourierValidatedAt, &orderValidation.CourierError, &orderValidation.UpdatedAt)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOrderValidationNotFound
	}

	return &orderValidation, nil
}

// SaveOrderValidation creates or updates order validation
func (repo *OrderRepository) SaveOrderValidation(
	ctx context.Context,
	orderValidation *domain.OrderValidation,
) error {
	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "INSERT INTO order_validations(order_id, courier_validated_at, courier_error) VALUES ($1, $2, $3)"
	_, err := client.ExecContext(
		ctx,
		query,
		orderValidation.OrderID,
		orderValidation.CourierValidatedAt,
		orderValidation.CourierError,
	)

	return err
}

// UpdateOrderValidation updates order validation when order validation was added
func (repo *OrderRepository) UpdateOrderValidation(
	ctx context.Context,
	orderValidation *domain.OrderValidation,
) error {

	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "UPDATE  order_validations SET courier_validated_at = $2, courier_error = $3, updated_at = $4 WHERE updated_at=$5 AND order_id=$1"
	result, err := client.ExecContext(
		ctx,
		query,
		orderValidation.OrderID,
		orderValidation.CourierValidatedAt,
		orderValidation.CourierError,
		time.Now(),
		orderValidation.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected > 0 {
		return nil
	}

	return domain.ErrOrderValidationNotFound
}

// GetOrderByID gets order status and order id from database by uuid order and return model with order id.
func (repo *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	var client Client

	if ctx.Value("transaction") != nil {
		client = ctx.Value("transaction").(*sql.Tx)
	} else {
		client = repo.client
	}

	query := "SELECT id, status, courier_id FROM orders WHERE id=$1 FOR SHARE"
	row := client.QueryRowContext(
		ctx,
		query,
		orderID,
	)

	var orderRow domain.Order
	err := row.Scan(&orderRow.ID, &orderRow.Status, &orderRow.CourierID)

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
