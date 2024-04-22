package domain_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	lm "github.com/Sharykhin/go-delivery-dymas/location/location_mocks"
)

// TestSaveLatestCourierLocation test three scenario success save and fail save in db and fail publish in third system
func TestSaveLatestCourierLocation(t *testing.T) {
	mc := minimock.NewController(t)
	c := qt.New(t)
	t.Run("success scenarios save latest geo position", func(t *testing.T) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad425",
			Latitude:  53.92680546122101,
			Longitude: 27.606307389240364,
			CreatedAt: time.Now(),
		}
		courierLocationRepositoryMock := lm.NewCourierLocationRepositoryInterfaceMock(mc)
		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			When(minimock.AnyContext, &courier).Then(nil)
		publisherLocationMock := lm.NewCourierLocationPublisherInterfaceMock(mc)
		publisherLocationMock.PublishLatestCourierLocationMock.
			When(minimock.AnyContext, &courier).Then(nil)
		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		err := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(err, qt.ErrorIs, nil)
	})
	t.Run("fail scenarios save latest geo position", func(t *testing.T) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad477",
			Latitude:  53.92,
			Longitude: 27.606,
			CreatedAt: time.Now(),
		}
		err := errors.New("repository error")
		courierLocationRepositoryMock := lm.NewCourierLocationRepositoryInterfaceMock(mc)
		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			When(minimock.AnyContext, &courier).Then(err)
		publisherLocationMock := lm.NewCourierLocationPublisherInterfaceMock(mc)
		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		err = fmt.Errorf("failed to store latest courier location in the repository: %w", err)
		errResult := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(err.Error(), qt.Equals, errResult.Error())
	})

	t.Run("fail scenarios publish latest geo position in third system", func(t *testing.T) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-data89ad477",
			Latitude:  53.42,
			Longitude: 27.106,
			CreatedAt: time.Now(),
		}
		courierLocationRepositoryMock := lm.NewCourierLocationRepositoryInterfaceMock(mc)
		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			When(minimock.AnyContext, &courier).Then(nil)
		publisherLocationMock := lm.NewCourierLocationPublisherInterfaceMock(mc)
		err := errors.New("publisher error")
		publisherLocationMock.PublishLatestCourierLocationMock.
			When(minimock.AnyContext, &courier).Then(err)
		err = fmt.Errorf("failed to publish latest courier location: %w", err)
		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		errResult := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(err.Error(), qt.Equals, errResult.Error())
	})
}
