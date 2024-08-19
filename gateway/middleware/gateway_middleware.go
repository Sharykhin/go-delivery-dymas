package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
	"github.com/gorilla/mux"
)

var uuid = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

// UuidMiddleware Check parameter by nane from request param is uuid or not and return Bad Request if param is not uuid
func UuidMiddleware(paramName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			param, ok := params[paramName]
			if !ok {
				message := fmt.Sprintf("%s doesn't exist in the defined route path", param)
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(&pkghttp.ResponseMessage{
					Status:  "Error",
					Message: message,
				})

				if err != nil {
					log.Printf("failed to encode json response: %v\n", err)
				}
				return
			}
			isUuid := uuid.MatchString(param)
			if !isUuid {
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(&pkghttp.ResponseMessage{
					Status:  "Error",
					Message: "Invalid UUID format",
				})

				if err != nil {
					log.Printf("failed to encode json response: %v\n", err)
				}

				return
			}

			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
}
