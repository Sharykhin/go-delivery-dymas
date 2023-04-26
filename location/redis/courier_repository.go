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

func (repo *CourierLatestGeoPositionRepository) SaveLatestCourierGeoPosition(courierID string, latitude, longitude float64) error {
	// add locations to the database
	err := repo.client.GeoAdd(repo.ctx, repo.indexGeo,
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

func CreateCouriersRepository(ctx context.Context, options *coreredis.Options, indexGeo string) *CourierLatestGeoPositionRepository {
	client, err := CreateConnect(ctx, options)
	if err != nil {
		panic(err)
	}

	return &CourierLatestGeoPositionRepository{
		indexGeo: indexGeo,
		client:   client,
		ctx:      ctx,
	}
}
