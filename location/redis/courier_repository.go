package redis

import (
	"context"
	"fmt"

	coreredis "github.com/redis/go-redis/v9"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

// CourierLocationRepository needs for managing location courier in Redis.
// At the current moment it provides API to store only the latest courier position.
type CourierLocationRepository struct {
	indexGeo string
	client   *coreredis.Client
}

const courierLatestCordsKey = "courier_latest_cord"

func (repo *CourierLocationRepository) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *domain.CourierLocation) error {
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

// NewCourierLocationRepository  creates courier location repository.
func NewCourierLocationRepository(client *coreredis.Client) *CourierLocationRepository {
	return &CourierLocationRepository{
		indexGeo: courierLatestCordsKey,
		client:   client,
	}
}
