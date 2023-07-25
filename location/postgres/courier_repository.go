package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	_ "github.com/lib/pq"
)

const table_name = "courier_latest_cord"
const user_db = "citizix_user"
const password_db = "S3cret"
const db_name = "courier_location"

type CourierLocationRepository struct {
	indexGeo string
	client   *sql.DB
}

func (repo *CourierLocationRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *domain.CourierLocation) error {
	query := fmt.Sprintf("insert into %s (courier_id, latitude, longitude, created_at) values ($1, $2, $3, $4)", table_name)
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courierLocation.CourierID,
		courierLocation.Latitude,
		courierLocation.Longitude,
		courierLocation.CreatedAt,
	)
	if row.Err() != nil {
		fmt.Println(row.Err())
		return fmt.Errorf("Row couirier location was not saved: %w", row.Err())
	}

	return nil
}

func NewCourierLocationRepository() (*CourierLocationRepository, error) {
	connPostgres := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user_db, password_db, db_name)
	db, err := sql.Open("postgres", connPostgres)
	if err != nil {
		db.Close()
		fmt.Println(err)
		return nil, err
	}
	courierLocationRepository := CourierLocationRepository{
		client: db,
	}

	return &courierLocationRepository, nil
}
