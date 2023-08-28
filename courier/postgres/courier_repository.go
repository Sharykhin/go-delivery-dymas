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

func NewCourierRepository(client *sql.DB) *CourierRepository {

	courierRepository := CourierRepository{
		client: client,
	}

	return &courierRepository
}
