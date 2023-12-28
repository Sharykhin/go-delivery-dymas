package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (repo *CourierRepository) GetAppliedCourier(ctx context.Context) (*domain.Courier, error) {

	query := "UPDATE courier SET is_available = FALSE " +
		"where id = (SELECT id FROM courier WHERE is_available = TRUE LIMIT 1 FOR UPDATE SKIP LOCK) RETURNING id"
	row := repo.client.QueryRowContext(
		ctx,
		query,
	)

	var courierRow domain.Courier

	err := row.Scan(&courierRow.ID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCourierNotFound
	}

	if err != nil {
		return nil, err
	}

	return &courierRow, nil
}

// NewCourierRepository creates new courier repository
func NewCourierRepository(client *sql.DB) *CourierRepository {
	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
