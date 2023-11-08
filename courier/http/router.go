package http

import (
	nethttp "net/http"

	"github.com/gorilla/mux"
)

const UuidRegexp = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"

// Route handles different path routes
type Route struct {
	Handler func(nethttp.ResponseWriter, *nethttp.Request)
	Method  string
}

// NewCourierRoute creates for handling different path routes
func NewCourierRoute(routes map[string]Route, router *mux.Router) *mux.Router {
	for url, route := range routes {
		router.HandleFunc(url, route.Handler).Methods(route.Method)
	}

	return router
}
