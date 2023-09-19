package postgres

import (
	"context"
	"database/sql"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	_ "github.com/lib/pq"
)

type CourierRepository struct {
	client *sql.DB
}

func (repo CourierRepository) SaveCourier(ctx context.Context, courier *domain.Courier) (*domain.Courier, error) {
	query := "insert into courier (first_name, is_available) values ($1,$2) RETURNING id, first_name, is_available"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courier.FirstName,
		courier.IsAvailable,
	)

	var courierRow domain.Courier
	err := row.Scan(&courierRow.Id, &courierRow.FirstName, &courierRow.IsAvailable)

	return &courierRow, err
}

func (repo CourierRepository) GetCourierById(ctx context.Context, courierID string) (*domain.Courier, error)  {
	query := "SELECT * FROM courier WHERE id=$1"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courierID,
	)

	var courierRow domain.Courier
	err := row.Scan(&courierRow.Id, &courierRow.FirstName, &courierRow.IsAvailable)

	return &courierRow, err
}

func NewCourierRepository(client *sql.DB) *CourierRepository {

	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
