package redis

import (
	"context"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	coreredis "github.com/redis/go-redis/v9"
)

type CourierRepository struct {
	indexGeo string
	client   *coreredis.Client
}

const courierLatestCordsKey = "courier_latest_cord"

func (repo *CourierRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *domain.CourierLocation) error {
	// add locations to the database
	err := repo.client.GeoAdd(ctx, repo.indexGeo,
		&coreredis.GeoLocation{
			Name:      courierLocation.CourierID,
			Latitude:  courierLocation.Latitude,
			Longitude: courierLocation.Longitude,
		}).Err()
	if err != nil {
		return fmt.Errorf("failed to add courier geo location into redis: %w", err)
	}

	return nil
}

func CreateCouriersRepository(client *coreredis.Client) *CourierRepository {

	return &CourierRepository{
		indexGeo: courierLatestCordsKey,
		client:   client,
	}
}
