package domain

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"

	lm "github.com/Sharykhin/go-delivery-dymas/location/location_mocks"
)

type TestKeySaveLatestCourierLocation struct {
	context             context.Context
	courier             CourierLocation
	resultRepository    error
	resultPublisher     error
	resultServices      error
	descriptionTestCase string
}

var testKeysSaveLatestCourierLocation = []TestKeySaveLatestCourierLocation{
	TestKeySaveLatestCourierLocation{
		context: minimock.AnyContext,
		courier: CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad425",
			Latitude:  53.92680546122101,
			Longitude: 27.606307389240364,
			CreatedAt: time.Now(),
		},
		resultRepository:    nil,
		resultPublisher:     nil,
		resultServices:      nil,
		descriptionTestCase: "success scenarios save latest geo position",
	}, TestKeySaveLatestCourierLocation{
		context: minimock.AnyContext,
		courier: CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-d5d6d89ad477",
			Latitude:  53.92,
			Longitude: 27.606,
			CreatedAt: time.Now(),
		},
		resultRepository:    errors.New("repository error"),
		resultServices:      fmt.Errorf("failed to store latest courier location in the repository: %w", errors.New("repository error")),
		resultPublisher:     nil,
		descriptionTestCase: "fail scenarios save latest geo position",
	},
	TestKeySaveLatestCourierLocation{
		context: minimock.AnyContext,
		courier: CourierLocation{
			CourierID: "23906828-0744-4a48-a2ca-data89ad477",
			Latitude:  53.42,
			Longitude: 27.106,
			CreatedAt: time.Now(),
		},
		resultRepository:    nil,
		resultPublisher:     errors.New("publisher error"),
		resultServices:      fmt.Errorf("failed to publish latest courier location: %w", errors.New("publisher error")),
		descriptionTestCase: "fail scenarios publish latest geo position in third system",
	},
}

func TestSaveLatestCourierLocation(t *testing.T) {
	mc := minimock.NewController(t)
	c := qt.New(t)
	for _, testKey := range testKeysSaveLatestCourierLocation {
		courierLocationRepositoryMock := lm.NewCourierLocationRepositoryInterfaceMock(mc)
		courierLocationRepositoryMock.SaveLatestCourierGeoPositionMock.
			When(minimock.AnyContext, &testKey.courier).Then(testKey.resultRepository)
		publisherLocationMock := lm.NewCourierLocationPublisherInterfaceMock(mc)
		publisherLocationMock.PublishLatestCourierLocationMock.
			When(minimock.AnyContext, &testKey.courier).Then(testKey.resultPublisher)
		courierLocationService := NewCourierLocationService(courierLocationRepositoryMock, publisherLocationMock)
		err := courierLocationService.SaveLatestCourierLocation(minimock.AnyContext, &testKey.courier)
		c.Assert(err, qt.ErrorAs, testKey.resultServices)
	}

}
