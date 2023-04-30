package redis

import (
	"context"
	"fmt"
	coreredis "github.com/redis/go-redis/v9"
)

type CourierRepository struct {
	indexGeo string
	client   *coreredis.Client
	ctx      context.Context
}

const courierLatestCordsKey = "courier_latest_cord"

func (repo *CourierRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierID string, latitude, longitude float64) error {
	// add locations to the database
	err := repo.client.GeoAdd(ctx, repo.indexGeo,
		&coreredis.GeoLocation{
			Name:      courierID,
			Latitude:  latitude,
			Longitude: longitude,
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
