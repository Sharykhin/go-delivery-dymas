package handler

import (
	"log"
	nethttp "net/http"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// LocationPayload imagine payload from http query
type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

// LocationHandler handles request depending on location courier
type LocationHandler struct {
	courierLocationService domain.CourierLocationServiceInterface
	httpHandler            pkghttp.HandlerInterface
}

// NewLocationHandler creates location handler
func NewLocationHandler(
	courierLocationService domain.CourierLocationServiceInterface,
	handler pkghttp.HandlerInterface,
) *LocationHandler {
	return &LocationHandler{
		courierLocationService: courierLocationService,
		httpHandler:            handler,
	}
}

// HandlerCouriersLocation handles request depending on location courier and validate query have valid payload and save data from payload in storage
func (h *LocationHandler) HandlerCouriersLocation(w nethttp.ResponseWriter, r *nethttp.Request) {
	var locationPayload LocationPayload

	if err := h.httpHandler.DecodePayloadFromJson(r, &locationPayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	if err := h.httpHandler.ValidatePayload(&locationPayload); err != nil {
		h.httpHandler.FailResponse(w, err)

		return
	}

	vars := mux.Vars(r)
	courierID := vars["courier_id"]
	ctx := r.Context()
	courierLocation := domain.NewCourierLocation(
		courierID,
		locationPayload.Latitude,
		locationPayload.Longitude,
	)

	err := h.courierLocationService.SaveLatestCourierLocation(ctx, courierLocation)

	if err != nil {
		log.Printf("failed to store latest courier position: %v", err)

		h.httpHandler.FailResponse(w, err)

		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}
