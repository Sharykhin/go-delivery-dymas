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

type CourierRepositoryData struct {
	CourierID string
	Latitude  float64
	Longitude float64
}

func (repo *CourierRepository) SaveLatestCourierGeoPosition(data *domain.CourierRepositoryData, ctx context.Context) error {
	// add locations to the database
	err := repo.client.GeoAdd(ctx, repo.indexGeo,
		&coreredis.GeoLocation{
			Name:      data.CourierID,
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		}).Err()
	if err != nil {
		return fmt.Errorf("failed to add courier geo location into redis: %w", err)
	}

	return nil
}

func CreateCouriersRepository(client *coreredis.Client) domain.CourierRepositoryInterface {

	return &CourierRepository{
		indexGeo: courierLatestCordsKey,
		client:   client,
	}
}
