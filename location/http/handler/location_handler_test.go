package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	locationHandler "github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/Sharykhin/go-delivery-dymas/location/mock"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// TestHandlerCouriersLocation coverages scenarios failed decode and failed validation and success save
func TestHandlerCouriersLocation(t *testing.T) {
	c := qt.New(t)
	mc := minimock.NewController(c)
	handler := pkghttp.NewHandler()
	c.Run("failed to decode payload", func(c *qt.C) {
		req := httptest.NewRequest(http.MethodPost, "/courier/77204924-4714-40cd-845e-36fcc67f9479/location", nil)

		workerPoolMock := mock.NewCourierLocationWorkerPoolMock(mc)

		w := httptest.NewRecorder()

		locationHandler := locationHandler.NewLocationHandler(workerPoolMock, handler)
		locationHandler.HandlerCouriersLocation(w, req)

		res := w.Result()

		defer res.Body.Close()

		c.Assert(res.StatusCode, qt.Equals, http.StatusBadRequest)
	})

	c.Run("failed payload validation", func(c *qt.C) {

		bodyReader := bytes.NewReader([]byte(`{"latitude": 0, "longitude": 0}`))

		req := httptest.NewRequest(http.MethodPost, "/courier/77204924-4714-40cd-845e-36fcc67f9479/location", bodyReader)

		workerPoolMock := mock.NewCourierLocationWorkerPoolMock(mc)

		w := httptest.NewRecorder()
		locationHandler := locationHandler.NewLocationHandler(workerPoolMock, handler)
		locationHandler.HandlerCouriersLocation(w, req)

		res := w.Result()

		defer res.Body.Close()

		c.Assert(res.StatusCode, qt.Equals, http.StatusBadRequest)
	})

	c.Run("success save courier location", func(c *qt.C) {
		courierLocationTest := domain.CourierLocation{
			Latitude:  20,
			Longitude: 131,
			CreatedAt: time.Now(),
		}

		bodyReader := bytes.NewReader([]byte(`{"latitude": 20, "longitude": 131}`))

		req := httptest.NewRequest(http.MethodPost, "/courier/77204924-4714-40cd-845e-36fcc67f9479/location", bodyReader)

		workerPoolMock := mock.NewCourierLocationWorkerPoolMock(mc)
		workerPoolMock.AddTaskMock.Set(func(courierLocation *domain.CourierLocation) {
			c.Assert(courierLocation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &courierLocationTest)
		})

		w := httptest.NewRecorder()

		locationHandler := locationHandler.NewLocationHandler(workerPoolMock, handler)

		locationHandler.HandlerCouriersLocation(w, req)

		res := w.Result()

		defer res.Body.Close()

		c.Assert(res.StatusCode, qt.Equals, http.StatusNoContent)
	})
}
