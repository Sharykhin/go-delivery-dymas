package handler

import (
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
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

func (h *LocationHandler) validatePayload(payload *LocationPayload) (isValid bool, response *ResponseMessage) {
	err := h.validate.Struct(payload)
	if err != nil {
		var errorMessage string

		for _, errStruct := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("Incorrect Value %s %f", errStruct.StructField(), errStruct.Value())
			errorMessage += message + ","
		}

		errorMessage = strings.Trim(errorMessage, ",")

		return false, &ResponseMessage{
			Status:  "Error",
			Message: errorMessage,
		}

	}

	return true, nil
}

// HandlerCouriersLocation gets latest courier position
func (h *LocationHandler) HandlerCouriersLocation(w nethttp.ResponseWriter, r *nethttp.Request) {
	var locationPayload LocationPayload
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&locationPayload)

	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: "Incorrect json! Please check your JSON formating.",
		})

		if err != nil {
			log.Printf("failed to encode json response error: %v\n", err)
		}

		return
	}

	if isValid, response := h.validatePayload(&locationPayload); !isValid {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}
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
		err := json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: "Server Error.",
		})
		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}
		w.WriteHeader(nethttp.StatusInternalServerError)

		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}
