package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	_ "github.com/lib/pq"
)

const user_db = "citizix_user"
const password_db = "S3cret"
const db_name = "courier_location"

type CourierLocationRepository struct {
	Client *sql.DB
}

func (repo *CourierLocationRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *domain.CourierLocation) error {
	query := "insert into courier_latest_cord (courier_id, latitude, longitude, created_at) values ($1, $2, $3, $4)"
	_, err := repo.Client.ExecContext(
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

func NewCourierLocationRepository(dbName string, dbUser string, dbPassword string) (*CourierLocationRepository, error) {
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connPostgres)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	courierLocationRepository := CourierLocationRepository{
		Client: db,
	}

	return &courierLocationRepository, nil
}
