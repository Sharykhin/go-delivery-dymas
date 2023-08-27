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

type CourierRepository struct {
	client *sql.DB
}

func (repo CourierRepository) SaveCourier(ctx context.Context, courier domain.CourierModel) error {
	query := "insert into couriers (id, first_name, is_available) values ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err := repo.client.ExecContext(
		ctx,
		query,
		courier.Id,
		courier.FirstName,
		courier.IsAvailable,
	)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Row couirier was not saved: %w", err)
	}

	return nil
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

func NewCourierLocationRepository(client *sql.DB) *CourierLocationRepository {

	courierLocationRepository := CourierLocationRepository{
		client: client,
	}

	return &courierLocationRepository
}

func NewCourierRepository(client *sql.DB) *CourierRepository {

	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
