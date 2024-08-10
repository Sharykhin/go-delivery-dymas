package http

import (
	"net/http"
	nethttp "net/http"

	"github.com/gorilla/mux"
)

// Route handles different path routes
type Route struct {
	Handler     func(nethttp.ResponseWriter, *nethttp.Request)
	Methods     []string
	Middlewares []func(next http.Handler) http.Handler
}

// NewRoute creates for handling different path routes
func NewRoute(routes map[string]Route, router *mux.Router) *mux.Router {
	for url, route := range routes {
		if route.Middlewares != nil {
			handle := prepareMiddleware(http.HandlerFunc(route.Handler), route.Middlewares...)
			router.Handle(url, handle).Methods(route.Methods...)
		} else {
			router.HandleFunc(url, route.Handler).Methods(route.Methods...)
		}
	}

	return router
}

func prepareMiddleware(handler http.Handler, Middleware ...func(next http.Handler) http.Handler) http.Handler {
	for i := len(Middleware); i > 0; i-- {
		handler = Middleware[i-1](handler)
	}
	return handler
}
