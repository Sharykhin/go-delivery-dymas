package redis

import (
	"context"
	coreredis "github.com/redis/go-redis/v9"
)

type CourierRepository struct {
	indexGeo string
	client   *coreredis.Client
	ctx      context.Context
}

func (repo *CourierRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierID string, latitude, longitude float64) error {
	// add locations to the database
	err := repo.client.GeoAdd(ctx, repo.indexGeo,
		&coreredis.GeoLocation{
			Name:      courierID,
			Latitude:  latitude,
			Longitude: longitude,
		}).Err()
	if err != nil {
		return err
	}

	return nil
}

func CreateCouriersRepository(client *coreredis.Client, indexGeo string) *CourierRepository {

	return &CourierRepository{
		indexGeo: indexGeo,
		client:   client,
	}
}