package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
)

// CourierRepository saves and gets courier in db
type CourierRepository struct {
	client *sql.DB
}

// SaveCourier saves courier in db
func (repo *CourierRepository) SaveCourier(ctx context.Context, courier *domain.Courier) (*domain.Courier, error) {
	query := "insert into courier (first_name, is_available) values ($1,$2) RETURNING id, first_name, is_available"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courier.FirstName,
		courier.IsAvailable,
	)

	var courierRow domain.Courier
	err := row.Scan(&courierRow.ID, &courierRow.FirstName, &courierRow.IsAvailable)

	return &courierRow, err
}

// GetCourierByID gets courier by id from db
func (repo *CourierRepository) GetCourierByID(ctx context.Context, courierID string) (*domain.Courier, error) {
	query := "SELECT id,first_name,is_available  FROM courier WHERE id=$1"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courierID,
	)

	var courierRow domain.Courier
	err := row.Scan(&courierRow.ID, &courierRow.FirstName, &courierRow.IsAvailable)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCourierNotFound
	}

	if err != nil {
		return nil, err
	}

	return &courierRow, err
}

// AssignOrderToCourier assigns a free courier to order. It runs a transaction and after finding an available courier it inserts a record into order_assignments table. In case of concurrent request and having a conflict it just does nothing and returns already assigned courier
func (repo *CourierRepository) AssignOrderToCourier(ctx context.Context, orderID string) (courierAssignments domain.CourierAssignments, err error) {
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

	query := "UPDATE courier SET is_available = FALSE " +
		"where id = (SELECT id FROM courier WHERE is_available = TRUE LIMIT 1 FOR UPDATE) RETURNING id"
	row := tx.QueryRowContext(
		ctx,
		query,
	)

	var courierID string

	err = row.Scan(&courierID)

	if errors.Is(err, sql.ErrNoRows) {
		return
	}

	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock($1)", courierAssignments.Hash(orderID))
	if err != nil {
		return
	}

	query = "SELECT courier_id, order_id, created_at FROM order_assignments WHERE order_id=$1"
	row = tx.QueryRowContext(
		ctx,
		query,
		courierID,
	)

	err = row.Scan(&courierAssignments.CourierID, &courierAssignments.OrderID, &courierAssignments.CreatedAt)
	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	query = "INSERT INTO order_assignments (order_id, courier_id, created_at) VALUES ($1, $2, $3)" +
		" RETURNING courier_id, order_id, created_at"

	row = tx.QueryRowContext(
		ctx,
		query,
		orderID,
		courierID,
		time.Now(),
	)

	err = row.Scan(&courierAssignments.CourierID, &courierAssignments.OrderID, &courierAssignments.CreatedAt)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

// NewCourierRepository creates new courier repository
func NewCourierRepository(client *sql.DB) *CourierRepository {
	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
