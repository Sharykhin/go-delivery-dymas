package handler

import (
	"log"
	nethttp "net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// CourierHandler handles courier request.
type CourierHandler struct {
	validate          *validator.Validate
	courierRepository domain.CourierRepositoryInterface
	courierService    *domain.CourierService
	httpHandler       pkghttp.Handler
}

// CourierPayload passes payload in courier create request.
type CourierPayload struct {
	FirstName string `json:"first_name" validate:"required"`
}

// NewCourierHandler  creates courier handler.
func NewCourierHandler(
	courierService *domain.CourierService,
) *CourierHandler {
	return &CourierHandler{
		validate:       validator.New(),
		courierService: courierService,
	}
}

// HandlerCourierCreate handles request create courier.
func (h *CourierHandler) HandlerCourierCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	if err := h.httpHandler.EncodeResponseToJson(w, courier); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	w.WriteHeader(nethttp.StatusCreated)
}

// GetCourier handles request get courier.
func (h *CourierHandler) GetCourier(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ctx := r.Context()
	courierID := vars["id"]
	courierResponse, err := h.courierService.GetCourierWithLatestPosition(ctx, courierID)

	if err != nil {
		log.Printf("failed to save courier: %v", err)
		h.httpHandler.FailResponse(w, err)

		return
	}

	if err := h.httpHandler.EncodeResponseToJson(w, courierResponse); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	w.WriteHeader(nethttp.StatusOK)
}
