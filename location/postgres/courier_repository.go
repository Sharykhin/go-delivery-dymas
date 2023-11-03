package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

type CourierLocationRepository struct {
	client *sql.DB
}

type ErrorTypeHandlerDatabase struct {
	errors []error
}

func (e *ErrorTypeHandlerDatabase) CompareTypeError(err error) bool {
	for _, errCompare := range e.errors {

		isErrorExists := errors.Is(errCompare, err)

		if isErrorExists {
			return true
		}
	}

	return false
}

func (e *ErrorTypeHandlerDatabase) AddTypeError(err error) {
	e.errors = append(e.errors, err)
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
		return fmt.Errorf("row couirier location was not saved: %w", err)
	}

	return nil
}

func (repo *CourierLocationRepository) GetLatestPositionCourierByID(ctx context.Context, courierID string) (*domain.CourierLocation, error) {
	query := "SELECT latitude, longitude  FROM courier_latest_cord WHERE courier_id=$1 ORDER BY created_at DESC LIMIT 1"
	row := repo.client.QueryRowContext(
		ctx,
		query,
		courierID,
	)

	var courierLocationRow domain.CourierLocation
	err := row.Scan(&courierLocationRow.Latitude, &courierLocationRow.Longitude)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCourierNotFound
	}

	return &courierLocationRow, err
}

func NewCourierLocationRepository(client *sql.DB) *CourierLocationRepository {
	courierLocationRepository := CourierLocationRepository{
		client: client,
	}

	return &courierLocationRepository
}

func NewErrorTypeHandleDatabase() *ErrorTypeHandlerDatabase {
	var errorTypeHandlerDatabase *ErrorTypeHandlerDatabase
	errorTypeHandlerDatabase.AddTypeError(sql.ErrConnDone)
	return errorTypeHandlerDatabase
}
