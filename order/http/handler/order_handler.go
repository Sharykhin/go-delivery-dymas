package handler

import (
	"errors"
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// OrderCreatePayload imagine payload from http query.
type OrderCreatePayload struct {
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
}

// OrderResponse imagine response order status from http query.
type OrderResponse struct {
	Status string `json:"status"`
	ID     string `json:"order_id"`
}

// OrderStatusResponse imagine response order status from http query by order id.
type OrderStatusResponse struct {
	Status string `json:"status"`
}

// OrderHandler handles courier request.
type OrderHandler struct {
	orderService domain.OrderService
	httpHandler  pkghttp.HandlerInterface
}

// NewOrderHandler creates order handler.
func NewOrderHandler(
	orderService domain.OrderService,
	handler pkghttp.HandlerInterface,
) *OrderHandler {
	return &OrderHandler{
		httpHandler:  handler,
		orderService: orderService,
	}
}

// HandleOrderCreate handles request order and validate query have valid payload and save data from payload in storage.
func (h *OrderHandler) HandleOrderCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	var orderCreatePayload OrderCreatePayload

	if err := h.httpHandler.DecodePayloadFromJson(r, &orderCreatePayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	if err := h.httpHandler.ValidatePayload(&orderCreatePayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	ctx := r.Context()

	order := domain.NewOrder(orderCreatePayload.PhoneNumber)
	order, err := h.orderService.CreateOrder(
		ctx,
		order,
	)

	if err != nil {
		log.Printf("Failed to save courier: %v", err)
		h.httpHandler.FailResponse(w, err)

		return
	}

	orderStatusResponse := OrderResponse{
		Status: order.Status,
		ID:     order.ID,
	}

	h.httpHandler.SuccessResponse(w, orderStatusResponse, nethttp.StatusAccepted)
}

// HandleGetByOrderID GetStatusByOrderId handles request and return order id and order status.
func (h *OrderHandler) HandleGetByOrderID(w nethttp.ResponseWriter, r *nethttp.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]
	order, err := h.orderService.GetOrderByID(r.Context(), orderID)

	if err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	orderStatusResponse := OrderStatusResponse{
		Status: order.Status,
	}

	h.httpHandler.SuccessResponse(w, orderStatusResponse, nethttp.StatusOK)
}

// HandleOrderCancel handles cancel order and publish event cancel in kafka for removing courier order assign.
func (h *OrderHandler) HandleOrderCancel(w nethttp.ResponseWriter, r *nethttp.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]
	err := h.orderService.CancelOrderByID(r.Context(), orderID)

	if errors.Is(err, domain.ErrorCanceledOrder) {
		err = fmt.Errorf("conflict cancel order: %w, %w", err, pkghttp.ErrConflict)
	}

	if err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	orderStatusResponse := OrderResponse{
		Status: domain.OrderStatusCanceled,
		ID:     orderID,
	}

	h.httpHandler.SuccessResponse(w, orderStatusResponse, nethttp.StatusOK)
}
