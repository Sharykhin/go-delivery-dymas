package handler

import (
	nethttp "net/http"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// LocationPayload imagine payload from http query.
type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

// LocationHandler handles request depending on location courier.
type LocationHandler struct {
	courierLocationWorkerPool domain.CourierLocationWorkerPool
	httpHandler               pkghttp.HandlerInterface
}

// NewLocationHandler creates location handler.
func NewLocationHandler(
	courierLocationWorkerPool domain.CourierLocationWorkerPool,
	handler pkghttp.HandlerInterface,
) *LocationHandler {
	return &LocationHandler{
		courierLocationWorkerPool: courierLocationWorkerPool,
		httpHandler:               handler,
	}
}

// HandlerCouriersLocation handles request depending on location courier and validate query have valid payload and save data from payload in storage.
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
	courierLocation := domain.NewCourierLocation(
		courierID,
		locationPayload.Latitude,
		locationPayload.Longitude,
	)

	h.courierLocationWorkerPool.AddTask(courierLocation)

	w.WriteHeader(nethttp.StatusNoContent)
}
