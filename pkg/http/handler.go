package http

import (
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ResponseMessage provides format response on courier request.
type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Handler struct {
	validate *validator.Validate
}

func (h *Handler) DecodePayloadFromJson(w nethttp.ResponseWriter, r *nethttp.Request, requestData any) bool {
	err := json.NewDecoder(r.Body).Decode(requestData)

	if err != nil {
		h.ErrorResponse("incorrect json! please check your json formatting.", w, nethttp.StatusBadRequest)

		return false
	}

	return true
}

func (h *Handler) EncodeResponseToJson(w nethttp.ResponseWriter, r *nethttp.Request, requestData any) bool {
	err := json.NewEncoder(w).Encode(requestData)

	if err != nil {
		h.ErrorResponse("failed to encode json response.", w, nethttp.StatusInternalServerError)

		return false
	}

	return true
}

func (h *Handler) ValidatePayload(w nethttp.ResponseWriter, payload any) bool {
	err := h.validate.Struct(payload)
	if err != nil {
		var errorMessage string

		for _, errStruct := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("Incorrect Value %s %f", errStruct.StructField(), errStruct.Value())
			errorMessage += message + ","
		}

		errorMessage = strings.Trim(errorMessage, ",")
		h.ErrorResponse(errorMessage, w, nethttp.StatusBadRequest)

		return false

	}

	return true
}

func (h *Handler) ErrorResponse(message string, w nethttp.ResponseWriter, codeStatus int) {
	err := json.NewEncoder(w).Encode(&ResponseMessage{
		Status:  "Error",
		Message: message,
	})

	if err != nil {
		log.Printf("failed to encode json response: %v\n", err)
	}

	w.WriteHeader(codeStatus)
}
