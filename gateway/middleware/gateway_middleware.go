package middleware

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// UuidMiddleware Check parameter by nane from request param is uuid or not and return Bad Request if param is not uuid
func UuidMiddleware(paramName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			queryParams := mux.Vars(r)
			param, ok := queryParams[paramName]
			isUUID, err := regexp.MatchString(
				"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
				param,
			)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !ok || !isUUID {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Println(r.RequestURI)

			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
}
