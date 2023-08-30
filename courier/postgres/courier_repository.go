package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	_ "github.com/lib/pq"
)

type CourierRepository struct {
	client *sql.DB
}

func (repo CourierRepository) SaveCourier(ctx context.Context, courier domain.Courier) error {
	query := "insert into courier (first_name) values ($1) ON CONFLICT DO NOTHING"
	_, err := repo.client.ExecContext(
		ctx,
		query,
		courier.FirstName,
	)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Row couirier was not saved: %w", err)
	}

	return nil
}

func NewCourierRepository(client *sql.DB) *CourierRepository {

	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
