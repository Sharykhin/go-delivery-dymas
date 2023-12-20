package handler

import (
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/order/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// OrderCreatePayload imagine payload from http query.
type OrderCreatePayload struct {
	Phone string `json:"phone_number" validate:"omitempty,e164"`
}

// OrderStatusResponse imagine response order status from http query.
type OrderStatusResponse struct {
	Status string `json:"status"`
	ID     string `json:"order_id"`
}

// OrderHandler handles courier request.
type OrderHandler struct {
	orderService *domain.OrderService
	httpHandler  pkghttp.HandlerInterface
}

// NewOrderHandler creates order handler.
func NewOrderHandler(
	orderService *domain.OrderService,
	handler pkghttp.HandlerInterface,
) *OrderHandler {
	return &OrderHandler{
		httpHandler:  handler,
		orderService: orderService,
	}
}

// HandlerOrderCreate handles request order and validate query have valid payload and save data from payload in storage.
func (h *OrderHandler) HandlerOrderCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
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

	order := domain.NewOrder(orderCreatePayload.Phone)
	order, err := h.orderService.CreateOrder(
		ctx,
		order,
	)

	orderStatusResponse := OrderStatusResponse{
		Status: order.Status,
		ID:     order.ID,
	}
	if err != nil {
		log.Printf("Failed to save courier: %v", err)
		h.httpHandler.FailResponse(w, err)

		return
	}

	h.httpHandler.SuccessResponse(w, orderStatusResponse, nethttp.StatusAccepted)
}

// HandlerOrderGetStatusByOrderId GetStatusByOrderId handles request and return order id and order status.
func (h *OrderHandler) HandlerOrderGetStatusByOrderId(w nethttp.ResponseWriter, r *nethttp.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]
	order, err := h.orderService.GetStatusByOrderId(r.Context(), orderID)
	orderStatusResponse := OrderStatusResponse{
		Status: order.Status,
		ID:     order.ID,
	}
	fmt.Println(orderStatusResponse)
	if err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}
	h.httpHandler.SuccessResponse(w, orderStatusResponse, nethttp.StatusOK)
}
