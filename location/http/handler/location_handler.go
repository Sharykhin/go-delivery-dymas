package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	nethttp "net/http"
	"strings"
)

type CourierRepository interface {
	SaveLatestCourierGeoPosition(ctx context.Context, courierID string, latitude, longitude float64) error
}

type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LocationHandler struct {
	validate          *validator.Validate
	courierRepository CourierRepository
}

func NewLocationHandler(courierRepository CourierRepository) *LocationHandler {
	return &LocationHandler{
		validate:          validator.New(),
		courierRepository: courierRepository,
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

func (h *LocationHandler) HandlerCouriersLocation(w nethttp.ResponseWriter, r *nethttp.Request) {
	var LocationPayload LocationPayload
	err := json.NewDecoder(r.Body).Decode(&LocationPayload)
	w.Header().Set("Content-Type", "application/json")

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

	if isValid, response := h.validatePayload(&LocationPayload); !isValid {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}
		return
	}
	vars := mux.Vars(r)
	courierId := vars["courier_id"]
	var ctx = context.Background()
	err = h.courierRepository.SaveLatestCourierGeoPosition(ctx, courierId, LocationPayload.Latitude, LocationPayload.Longitude)
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
