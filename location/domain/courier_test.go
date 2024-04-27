package domain_test

import (
	"errors"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/Sharykhin/go-delivery-dymas/location/mock"
)

// TestSaveLatestCourierLocation test three scenario success save and fail save in db and fail publish in third system
func TestSaveLatestCourierLocation(t *testing.T) {
	c := qt.New(t)
	mc := minimock.NewController(c)

	c.Run("success scenarios save latest geo position", func(c *qt.C) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad425",
			Latitude:  53.92680546122101,
			Longitude: 27.606307389240364,
			CreatedAt: time.Now(),
		}

		courierLocationRepositoryMock := mock.NewCourierLocationRepositoryInterfaceMock(mc)

		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			Expect(minimock.AnyContext, &courier).Return(nil)

		publisherLocationMock := mock.NewCourierLocationPublisherInterfaceMock(mc)

		publisherLocationMock.PublishLatestCourierLocationMock.
			Expect(minimock.AnyContext, &courier).Return(nil)

		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		err := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(err, qt.IsNil)
	})

	c.Run("fail scenarios save latest geo position", func(c *qt.C) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad477",
			Latitude:  53.92,
			Longitude: 27.606,
			CreatedAt: time.Now(),
		}

		courierLocationRepositoryMock := mock.NewCourierLocationRepositoryInterfaceMock(mc)

		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			Expect(minimock.AnyContext, &courier).Return(errors.New("repository error"))

		publisherLocationMock := mock.NewCourierLocationPublisherInterfaceMock(mc)

		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)

		errResult := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(
			errResult,
			qt.ErrorMatches,
			"failed to store latest courier location in the repository: repository error",
		)
	})

	c.Run("fail scenarios publish latest geo position in third system", func(c *qt.C) {
		courier := domain.CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-data89ad477",
			Latitude:  53.42,
			Longitude: 27.106,
			CreatedAt: time.Now(),
		}
		courierLocationRepositoryMock := mock.NewCourierLocationRepositoryInterfaceMock(mc)

		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			Expect(minimock.AnyContext, &courier).Return(nil)

		publisherLocationMock := mock.NewCourierLocationPublisherInterfaceMock(mc)

		publisherLocationMock.PublishLatestCourierLocationMock.
			Expect(minimock.AnyContext, &courier).Return(errors.New("publisher error"))

		courierLocationService := domain.NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		errResult := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &courier)
		c.Assert(errResult, qt.ErrorMatches, "failed to publish latest courier location: publisher error")
	})
}
