package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	_ "github.com/lib/pq"
)

type CourierLocationRepository struct {
	client *sql.DB
}

func (repo *CourierLocationRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *domain.CourierLocation) error {
	query := "insert into courier_latest_cord (courier_id, latitude, longitude, created_at) values ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	_, err := repo.client.ExecContext(
		ctx,
		query,
		courierLocation.CourierID,
		courierLocation.Latitude,
		courierLocation.Longitude,
		courierLocation.CreatedAt,
	)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Row couirier location was not saved: %w", err)
	}

	return nil
}

func (repo *CourierLocationRepository) GetLatestPositionCourierById(ctx context.Context, courierID string) (*domain.CourierLocation, error)  {
	query := "SELECT * FROM courier_latest_cord WHERE courier_id=$1"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courierID,
	)

	var courierLocationRow domain.CourierLocation
	err := row.Scan(&courierLocationRow.Longitude, &courierLocationRow.Latitude)

	return &courierLocationRow, err
}

func NewCourierLocationRepository(client *sql.DB) *CourierLocationRepository {

	courierLocationRepository := CourierLocationRepository{
		client: client,
	}

	return &courierLocationRepository
}
