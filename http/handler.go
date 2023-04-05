package courierService

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	nethttp "net/http"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var validate *validator.Validate

func handlerCouriersLocation(response http.ResponseWriter, request *http.Request) {
	var Location Location
	var errors [2]ResponseMessage
	err := json.NewDecoder(request.Body).Decode(&Location)
	response.Header().Set("Content-Type", "application/json")

	if err != nil {
		response.WriteHeader(nethttp.StatusBadRequest)
		messageError, errJson := json.Marshal(&ResponseMessage{
			Status:  "Error",
			Message: "Incorrect json! Please check your JSON formating.",
		})

		fmt.Println(messageError)
		if errJson == nil {
			fmt.Fprintln(response, string(messageError))
		}

		return
	}

	validate = validator.New()
	err = validate.Struct(&Location)
	if errStruct != nil {
		response.WriteHeader(400)

		if _, ok := errStruct.(*validator.InvalidValidationError); ok {
			fmt.Println(errStruct)
			fmt.Printf("Success2")
			return
		}

		increment := 0
		for _, errStruct := range errStruct.(validator.ValidationErrors) {
			message := fmt.Sprintf("Incorrect Value %s %f", errStruct.StructField(), errStruct.Value())
			errors[increment] = ResponseMessage{
				Status:  "Error",
				Message: message,
			}

			increment++
		}

		if len(errors) > 0 {
			messageError, errJson := json.Marshal(errors)
			if errJson == nil {
				fmt.Fprintln(response, string(messageError))
			}

		}

		return
	}

	response.WriteHeader(204)
}
