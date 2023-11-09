package handler

import (
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
)

// CourierHandler handles courier request.
type CourierHandler struct {
	validate          *validator.Validate
	courierRepository domain.CourierRepositoryInterface
	courierService    *domain.CourierService
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
	var courierPayload CourierPayload

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&courierPayload)

	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		err = json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: "Incorrect json! Please check your JSON formatting.",
		})

		if err != nil {
			log.Printf("failed to encode json response error: %v\n", err)
		}

		return
	}

	if isValid, response := h.validatePayload(&courierPayload); !isValid {
		w.WriteHeader(nethttp.StatusBadRequest)

		err = json.NewEncoder(w).Encode(response)

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

// GetCourier handles request get courier.
func (h *CourierHandler) GetCourier(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ctx := r.Context()
	courierID := vars["id"]
	courierResponse, err := h.courierService.GetCourierWithLatestPosition(ctx, courierID)

	if err != nil {
		h.errorHandler("Failed to get courier: %v", err, w, nethttp.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(courierResponse)
	if err != nil {
		h.errorHandler("Failed to encode json response: %v\n", err, w, nethttp.StatusInternalServerError)

		return
	}

	w.WriteHeader(nethttp.StatusOK)
}

func (h *CourierHandler) validatePayload(s any) (bool, *ResponseMessage) {
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
