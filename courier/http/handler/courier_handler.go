package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/courier/domain"
	pb "github.com/Sharykhin/go-delivery-dymas/proto/generate/location/v1"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	nethttp "net/http"
	"strings"
)

type CourierHandler struct {
	validate          *validator.Validate
	courierRepository domain.CourierRepositoryInterface
	courierClient     pb.CourierClient
}

type CourierPayload struct {
	FirstName string `json:"first_name" validate:"required"`
}

type LastPositionPayload struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type CourierLatestPositionPayload struct {
	CourierLatestPosition LastPositionPayload `json:"last_position"`
	FirstName             string              `json:"first_name" validate:"required"`
	Id                    string              `json:"id" validate:"uuid,required"`
	IsAvailable           bool                `json:"is_available" validate:"boolean,required"`
}
type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewCourierHandler(
	repo domain.CourierRepositoryInterface,
	courierClient pb.CourierClient,
) *CourierHandler {
	return &CourierHandler{
		validate:          validator.New(),
		courierRepository: repo,
		courierClient:     courierClient,
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
		h.internalServerErrorPrepare("Failed to save courier: %v", err, w)

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
	var courierLatestPositionPayload CourierLatestPositionPayload
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	ctx := r.Context()
	courierID := vars["id"]
	courier, err := h.courierRepository.GetCourierById(
		ctx,
		courierID,
	)
	if err != nil {
		h.internalServerErrorPrepare("Failed to get courier: %v", err, w)

		return
	}
	courierLatestPositionResponse, err := h.courierClient.GetCourierLatestPosition(ctx, &pb.GetCourierLatestPositionRequest{
		CourierId: courierID,
	})

	if err != nil {
		h.internalServerErrorPrepare("Failed to get last position courier: %v", err, w)

		return
	}
	courierLatestPosition := LastPositionPayload{
		Latitude:  courierLatestPositionResponse.Latitude,
		Longitude: courierLatestPositionResponse.Longitude,
	}

	courierLatestPositionPayload = CourierLatestPositionPayload{
		FirstName:             courier.FirstName,
		Id:                    courier.Id,
		IsAvailable:           courier.IsAvailable,
		CourierLatestPosition: courierLatestPosition,
	}

	err = json.NewEncoder(w).Encode(courierLatestPositionPayload)
	if err != nil {
		h.internalServerErrorPrepare("Failed to encode json response: %v\n", err, w)

		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (h *CourierHandler) validatePayload(payload *CourierPayload) (isValid bool, response *ResponseMessage) {
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

func (h *CourierHandler) internalServerErrorPrepare(message string, err error, w nethttp.ResponseWriter) {
	log.Printf(message, err)
	err = json.NewEncoder(w).Encode(&ResponseMessage{
		Status:  "Error",
		Message: "Server Error.",
	})
	if err != nil {
		log.Printf("failed to encode json response: %v\n", err)
	}

	w.WriteHeader(nethttp.StatusInternalServerError)
}
