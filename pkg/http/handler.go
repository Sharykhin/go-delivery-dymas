package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	nethttp "net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrDecodeFailed we return this error when we can not decode payload from http query
var ErrDecodeFailed = errors.New("failed to decode payload")

var ErrValidatePayloadFailed = errors.New("failed to validated payload")

// ResponseMessage returns when we have bad request, or we have problem on server
type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Handler abstract handler we can reuse it in different handlers
type Handler struct {
	Validator *validator.Validate
}

// DecodePayloadFromJson decodes payload from body http query and handle exceptions scenarios
func (h *Handler) DecodePayloadFromJson(r *nethttp.Request, requestData any) error {
	err := json.NewDecoder(r.Body).Decode(requestData)

	if err != nil {
		log.Printf("incorrect json! please check your json formatting: %v\n", err)

		return ErrDecodeFailed
	}

	return nil
}

// EncodeResponseToJson  Encodes response,that return user for http query and handle exceptions scenarios
func (h *Handler) EncodeResponseToJson(w nethttp.ResponseWriter, requestData any) error {
	err := json.NewEncoder(w).Encode(requestData)

	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		log.Panicf("failed to encode json response: %v\n", err)

	}

	return nil
}

// ValidatePayload validates some payload from http query
func (h *Handler) ValidatePayload(payload any) error {
	err := h.Validator.Struct(payload)

	if err != nil {

		return fmt.Errorf("%v:%w", err, ErrValidatePayloadFailed)

		return ErrValidatePayloadFailed
	}

	return nil
}

// FailResponse returns response for bad request
func (h *Handler) FailResponse(w nethttp.ResponseWriter, errFailResponse error) {

	if errors.Is(errFailResponse, ErrDecodeFailed) {

		err := json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: errFailResponse.Error(),
		})

		if err != nil {
			log.Printf("failed to encode json response: %v\n", err)
		}

		w.WriteHeader(nethttp.StatusBadRequest)

		return
	} else if errors.Is(errFailResponse, ErrValidatePayloadFailed) {
		var errorMessage string

		for _, errStruct := range ErrValidatePayloadFailed.(validator.ValidationErrors) {
			message := fmt.Sprintf("Incorrect Value %s %f", errStruct.StructField(), errStruct.Value())
			errorMessage += message + ","
		}

		errorMessage = strings.Trim(errorMessage, ",")
		json.NewEncoder(w).Encode(&ResponseMessage{
			Status:  "Error",
			Message: errorMessage,
		})

		w.WriteHeader(nethttp.StatusBadRequest)

		return

	}

	log.Printf("Server error: %v\n", errFailResponse)
	w.WriteHeader(nethttp.StatusInternalServerError)
}
