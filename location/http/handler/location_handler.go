package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	nethttp "net/http"
	"strings"
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
	validate *validator.Validate
}

var validate *validator.Validate

func (h *LocationHandler) HandlerCouriersLocation(w nethttp.ResponseWriter, r *nethttp.Request) {
	var LocationPayload LocationPayload
	var errorMessage string
	err := json.NewDecoder(r.Body).Decode(&LocationPayload)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		err := json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: "Incorrect json! Please check your JSON formating.",
		})

		if err != nil {
			fmt.Println(err)
		}

		return
	}
	err = h.validate.Struct(&LocationPayload)
	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)

		for _, errStruct := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("Incorrect Value %s %f", errStruct.StructField(), errStruct.Value())
			errorMessage += message + ","
		}

		if len(errorMessage) > 0 {
			errorMessage = strings.Trim(errorMessage, ",")
			err := json.NewEncoder(w).Encode(&ResponseMessage{
				Status:  "Error",
				Message: errorMessage,
			})
			if err != nil {
				fmt.Println(err)
			}

		}

		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func NewLocationHandler() *LocationHandler {
	return &LocationHandler{
		validate: validator.New(),
	}
}
