package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	nethttp "net/http"
	"strings"
)

type CourierHandler struct {
	validate          *validator.Validate
	courierRepository domain.CourierRepositoryInterface
	courierService    *domain.CourierService
}

type CourierPayload struct {
	FirstName string `json:"first_name" validate:"required"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewCourierHandler(
	CourierService *domain.CourierService,
) *CourierHandler {
	return &CourierHandler{
		validate:       validator.New(),
		courierService: CourierService,
	}
}

func (h *CourierHandler) HandlerCourierCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	var courierPayload CourierPayload
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&courierPayload)

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

	if isValid, response := h.validatePayload(&courierPayload); !isValid {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}
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
		h.errorHandler("Failed to save courier: %v", err, w, nethttp.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(courier)
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)

		log.Printf("failed to encode json response: %v\n", err)
		return
	}
	w.WriteHeader(nethttp.StatusCreated)
}

func (h *CourierHandler) HandlerGetCourierLatestPosition(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	ctx := r.Context()
	courierID := vars["id"]
	courierResponse, err := h.courierService.GetCourierWithLatestPosition(ctx, courierID)
	isErrorNotFound := err != nil && errors.Is(err, domain.ErrorNotFound)
	if isErrorNotFound {
		h.errorHandler("Failed to get last position courier: %v", err, w, nethttp.StatusNotFound)

		return
	}
	if err != nil && !isErrorNotFound {
		h.errorHandler("Failed to get last position courier: %v", err, w, nethttp.StatusInternalServerError)

		return
	}

	if isValid, response := h.validatePayload(&courierResponse); !isValid {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(courierResponse)
	if err != nil {
		h.errorHandler("Failed to encode json response: %v\n", err, w, nethttp.StatusInternalServerError)

		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (h *CourierHandler) validatePayload(s interface{}) (isValid bool, response *ResponseMessage) {
	err := h.validate.Struct(s)
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

func (h *CourierHandler) errorHandler(message string, err error, w nethttp.ResponseWriter, codeStatus int) {
	log.Printf(message, err)
	err = json.NewEncoder(w).Encode(&ResponseMessage{
		Status:  "Error",
		Message: "Server Error.",
	})
	if err != nil {
		log.Printf("failed to encode json response: %v\n", err)
	}

	w.WriteHeader(codeStatus)
}
