package handler

import (
	"log"
	nethttp "net/http"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// CourierHandler handles courier request.
type CourierHandler struct {
	courierRepository domain.CourierRepositoryInterface
	courierService    *domain.CourierService
	httpHandler       pkghttp.HandlerInterface
}

// CourierPayload passes payload in courier create request.
type CourierPayload struct {
	FirstName string `json:"first_name" validate:"required"`
}

// NewCourierHandler  creates courier handler.
func NewCourierHandler(
	courierService *domain.CourierService,
	handler pkghttp.HandlerInterface,
) *CourierHandler {
	return &CourierHandler{
		courierService: courierService,
		httpHandler:    handler,
	}
}

// HandlerCourierCreate handles request create courier.
func (h *CourierHandler) HandlerCourierCreate(w nethttp.ResponseWriter, r *nethttp.Request) {

	var courierPayload CourierPayload

	if err := h.httpHandler.DecodePayloadFromJson(r, &courierPayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	if err := h.httpHandler.ValidatePayload(&courierPayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	ctx := r.Context()
	courier, err := h.courierRepository.SaveCourier(
		ctx,
		&domain.Courier{
			FirstName:   courierPayload.FirstName,
			IsAvailable: true,
		},
	)

	if err != nil {
		log.Printf("Failed to save courier: %v", err)
		h.httpHandler.FailResponse(w, err)

		return
	}

	h.httpHandler.SuccessResponse(w, courier, nethttp.StatusCreated)

	w.WriteHeader(nethttp.StatusCreated)
}

// GetCourier handles request get courier.
func (h *CourierHandler) GetCourier(w nethttp.ResponseWriter, r *nethttp.Request) {

	vars := mux.Vars(r)
	ctx := r.Context()
	courierID := vars["id"]
	courierResponse, err := h.courierService.GetCourierWithLatestPosition(ctx, courierID)

	if err != nil {
		log.Printf("failed to get courier: %v", err)
		h.httpHandler.FailResponse(w, err)

		return
	}

	if err := h.httpHandler.SuccessResponse(w, courierResponse, nethttp.StatusOK); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}
}
