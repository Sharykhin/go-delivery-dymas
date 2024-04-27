package domain_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gojuno/minimock/v3"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	om "github.com/Sharykhin/go-delivery-dymas/order/order_mocks"
)

// TestValidateOrderForService covered scenarios fail save order validation and fail save order and fail publish in third system also success update order flow
func TestValidateOrderForService(t *testing.T) {
	c := qt.New(t)
	mc := minimock.NewController(c)

	c.Run("fail get order from db", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad425"

		err := errors.New("test get fail order")
		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(nil, err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad425"}`)
		err = fmt.Errorf("failed to get order: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail get order validation with type not expected error", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad555"

		err := errors.New("test get fail order validation")
		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).Return(nil, err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to get order validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail unmarshal payload", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad555"

		err := errors.New("invalid character 'o' looking for beginning of value")
		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidation := domain.OrderValidation{
			OrderID: orderID,
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidation, nil)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`ooooooo`)
		err = fmt.Errorf("failed to unmarshal courier payload: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail unmarshal payload", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad555"

		err := errors.New("invalid character 'o' looking for beginning of value")
		order := domain.Order{
			ID: orderID,
		}

		orderRepositoryMock.GetOrderByIDMock.Expect(minimock.AnyContext, orderID).Return(&order, nil)

		orderValidation := domain.OrderValidation{
			OrderID: orderID,
		}

		orderRepositoryMock.GetOrderValidationByIDMock.Expect(minimock.AnyContext, orderID).
			Return(&orderValidation, nil)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`ooooooo`)
		err = fmt.Errorf("failed to unmarshal courier payload: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail save order validation", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail save order validation")
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

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to save order in database during validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail order update validation", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail update order validation")
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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			fmt.Println(orderValidation)
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to save order in database during validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail order update validation", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail save  order validation")
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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to save order in database during validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail order update", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail save order")
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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(nil)

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to order order in database during validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("fail order update", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail save order")
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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(nil)

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(err)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to order order in database during validation: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("failed to publish a order", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

		orderID := "23906828-0744-4a48-a2ca-d5d6d89ad666"

		err := errors.New("fail save order")
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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(nil)

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(nil)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		orderPublisherMock.PublishOrderMock.Expect(minimock.AnyContext, &order, domain.EventOrderUpdated).Return(err)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		err = fmt.Errorf("failed to publish a order in the kafka: %w", err)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.ErrorMatches, err.Error())
	})

	c.Run("success update order", func(c *qt.C) {
		orderRepositoryMock := om.NewOrderRepositoryMock(mc)

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

		orderRepositoryMock.UpdateOrderValidationMock.Inspect(func(ctx context.Context, orderValidation *domain.OrderValidation) {
			c.Assert(orderValidation, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), &orderValidationTest)
		})

		orderRepositoryMock.UpdateOrderValidationMock.Return(nil)

		orderRepositoryMock.UpdateOrderMock.Expect(minimock.AnyContext, &order).Return(nil)

		orderPublisherMock := om.NewOrderPublisherMock(mc)

		orderPublisherMock.PublishOrderMock.Expect(minimock.AnyContext, &order, domain.EventOrderUpdated).Return(nil)

		serviceName := "courier"
		orderServiceManager := domain.NewOrderServiceManager(orderRepositoryMock, orderPublisherMock)

		validationInfo := []byte(`{"courier_id": "23906828-0744-4a48-a2ca-d5d6d89ad777"}`)
		errResult := orderServiceManager.ValidateOrderForService(minimock.AnyContext, serviceName, orderID, validationInfo)

		c.Assert(errResult, qt.IsNil)
	})
}
