package handler

import (
	"encoding/json"
	"log"
	nethttp "net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LocationHandler struct {
	validate               *validator.Validate
	courierLocationService domain.CourierLocationServiceInterface
	httpHandler            pkghttp.Handler
}

// NewLocationHandler creates location handler
func NewLocationHandler(
	courierLocationService domain.CourierLocationServiceInterface,
) *LocationHandler {
	return &LocationHandler{
		validate:               validator.New(),
		courierLocationService: courierLocationService,
	}
}

// HandlerCouriersLocation gets latest courier position
func (h *LocationHandler) HandlerCouriersLocation(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	var locationPayload LocationPayload
	err := json.NewDecoder(r.Body).Decode(&locationPayload)

	if isDecode := httpHandler.DecodePayloadFromJson(w, r, &locationPayload); !isDecode {
		log.Printf("failed to encode json response error: %v\n", err)
		return
	}

	if isValid := h.ValidatePayload(&locationPayload); !isValid {
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
	err = h.courierLocationService.SaveLatestCourierLocation(ctx, courierLocation)

	if err != nil {
		log.Printf("failed to store latest courier position: %v", err)

		httpHandler.ErrorResponse("Server Error.", w, nethttp.StatusInternalServerError)

		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}
