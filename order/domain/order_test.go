package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	"github.com/Sharykhin/go-delivery-dymas/order/mock"
)

// TestValidateOrderForService covered scenarios fail save order validation and fail save order and fail publish in third system also success update order flow
func TestValidateOrderForService(t *testing.T) {
	c := qt.New(t)
	mc := minimock.NewController(c)

	c.Run("fail get order from db", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad425"

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).
			Return(nil, errors.New("fail get order from db"))

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad425"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, "failed to get order: fail get order from db")
	})

	c.Run("fail get order validation with type not expected error", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad555"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderRepositoryMock.GetOrderValidationByIDMock.
			Expect(minimock.AnyContext, orderID).Return(nil, errors.New("test get fail order validation"))

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, "failed to get order validation: test get fail order validation")
	})

	c.Run("fail unmarshal payload", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad555"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidation := domain.OrderValidation{
			OrderID: orderID,
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidation, nil)

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`ooooooo`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(
			errResult,
			qt.ErrorMatches,
			"failed to unmarshal courier payload: invalid character 'o' looking for beginning of value",
		)
	})

	c.Run("fail save order validation", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(nil, domain.ErrOrderValidationNotFound)
		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.SaveOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) error {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return errors.New("fail save order validation")
		})

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(
			errResult,
			qt.ErrorMatches,
			"failed to save order in database during validation: fail save order validation",
		)
	})

	c.Run("fail order update validation", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidationTest, nil)

		orderRepositoryMock.UpdateOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) error {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return errors.New("fail update order validation")
		})

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(
			errResult,
			qt.ErrorMatches,
			"failed to save order in database during validation: fail update order validation",
		)
	})

	c.Run("fail order update", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidationTest, nil)

		orderRepositoryMock.UpdateOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) (err error) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return nil
		})

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(errors.New("fail order update"))

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(
			errResult,
			qt.ErrorMatches,
			"failed to order order in database during validation: fail order update",
		)
	})

	c.Run("fail order update", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidationTest, nil)

		orderRepositoryMock.UpdateOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) (err error) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return nil
		})

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(errors.New("fail update order"))

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, "failed to order order in database during validation: fail update order")
	})

	c.Run("failed to publish a order", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidationTest, nil)

		orderRepositoryMock.UpdateOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) (err error) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return nil
		})

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(nil)

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		orderPublisherMock.PublishOrderMock.Expect(minimock.AnyContext, &order, domain.EventOrderUpdated).
			Return(errors.New("fail save order"))

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)

		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, "failed to publish a order in the kafka: fail save order")
	})

	c.Run("success update order", func(c *qt.C) {
		orderRepositoryMock := mock.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidationTest := domain.OrderValidation{
			OrderID:            orderID,
			CourierValidatedAt: time.Now(),
			CourierError:       "",
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidationTest, nil)

		orderRepositoryMock.UpdateOrderValidationMock.Set(func(ctx context.Context, orderValidation *domain.OrderValidation) (err error) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)

			return nil
		})

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(nil)

		orderPublisherMock := mock.NewOrderPublisherMock(mc)

		orderPublisherMock.PublishOrderMock.Expect(minimock.AnyContext, &order, domain.EventOrderUpdated).Return(nil)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.IsNil)
	})
}
