package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/go-playground/validator/v10"
	"log"
	nethttp "net/http"
	"strings"
)

type CourierCreateHandler struct {
	validate          *validator.Validate
	courierRepository domain.CourierRepositoryInterface
}

func NewCourierCreateHandler(
	repo domain.CourierRepositoryInterface,
) *CourierCreateHandler {
	return &CourierCreateHandler{
		validate:          validator.New(),
		courierRepository: repo,
	}
}

func (h *CourierCreateHandler) HandlerCourierCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	var courierPayload domain.CourierModel
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
	err = h.courierRepository.SaveCourier(
		ctx,
		courierPayload,
	)
	if err != nil {
		log.Printf("failed to save latest courier: %v", err)
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

func (h *CourierCreateHandler) validatePayload(payload *domain.CourierModel) (isValid bool, response *ResponseMessage) {
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
