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

// ResponseMessage provides format response on courier request.
type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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

	if isDecode := httpHandler.DecodePayloadFromJson(w, r, &courierPayload); !isDecode {
		return
	}

	if isValid := httpHandler.ValidatePayload(w, &courierPayload); !isValid {
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
		httpHandler.ErrorResponse("Failed to save courier", w, nethttp.StatusInternalServerError)

		return
	}

	if isEncode := httpHandler.EncodeResponseToJson(w, r, courier); !isEncode {
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
		httpHandler.ErrorResponse("Failed to get courier", w, nethttp.StatusInternalServerError)

		return
	}

	if isEncode := httpHandler.EncodeResponseToJson(w, r, courierResponse); !isEncode {
		return
	}

	w.WriteHeader(nethttp.StatusOK)
}
